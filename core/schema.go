package core

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"cloud.google.com/go/bigquery"
)

// Schema represents a BigQuery schema.
type Schema struct {
	Fields []*FieldSchema
}

// FieldSchema represents a field in a BigQuery schema.
type FieldSchema struct {
	Name        string
	Type        bigquery.FieldType
	Description string
	Repeated    bool
	Required    bool
	Fields      []*FieldSchema
}

// InferSchema infers the BigQuery schema from a given JSON object.
func InferSchema(jsonData []byte) (*Schema, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	schema := &Schema{}
	for key, value := range data {
		fieldSchema, err := inferFieldSchema(key, value)
		if err != nil {
			return nil, err
		}
		if fieldSchema != nil {
			schema.Fields = append(schema.Fields, fieldSchema)
		}
	}

	sort.Slice(schema.Fields, func(i, j int) bool {
		return schema.Fields[i].Name < schema.Fields[j].Name
	})

	return schema, nil
}

func inferFieldSchema(name string, value interface{}) (*FieldSchema, error) {
	switch value := value.(type) {
	case string:
		return &FieldSchema{Name: name, Type: bigquery.StringFieldType, Required: true}, nil
	case float64:
		if float64(int64(value)) == value {
			return &FieldSchema{Name: name, Type: bigquery.IntegerFieldType, Required: true}, nil
		}
		return &FieldSchema{Name: name, Type: bigquery.FloatFieldType, Required: true}, nil
	case bool:
		return &FieldSchema{Name: name, Type: bigquery.BooleanFieldType, Required: true}, nil
	case map[string]interface{}:
		subSchema, err := InferSchema(mustMarshal(value))
		if err != nil {
			return nil, err
		}
		return &FieldSchema{Name: name, Type: bigquery.RecordFieldType, Fields: subSchema.Fields, Required: true}, nil
	case []interface{}:
		if len(value) == 0 {
			return nil, fmt.Errorf("cannot infer type from empty slice")
		}

		// Check the type of the first element to differentiate between primitive and record slices.
		if _, ok := value[0].(map[string]interface{}); ok {
			var mergedSchema *Schema
			for _, item := range value {
				itemSchema, err := InferSchema(mustMarshal(item))
				if err != nil {
					return nil, err
				}
				if mergedSchema == nil {
					mergedSchema = itemSchema
				} else {
					mergedSchema = MergeSchemas(mergedSchema, itemSchema)
				}
			}
			return &FieldSchema{Name: name, Type: bigquery.RecordFieldType, Fields: mergedSchema.Fields, Repeated: true, Required: true}, nil
		}
		
		// It's a slice of primitives.
		elemField, err := inferFieldSchema(name, value[0])
		if err != nil {
			return nil, err
		}
		elemField.Repeated = true
		return elemField, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", value)
	}
}

func MergeSchemas(s1, s2 *Schema) *Schema {
	fieldMap := make(map[string]*FieldSchema)

	for _, f := range s1.Fields {
		clone := *f
		fieldMap[f.Name] = &clone
	}

	for _, f2 := range s2.Fields {
		if f1, ok := fieldMap[f2.Name]; ok {
			if f1.Type != f2.Type {
				panic(fmt.Sprintf("type mismatch for field %s: %s vs %s", f1.Name, f1.Type, f2.Type))
			}
			if f1.Repeated != f2.Repeated {
				panic(fmt.Sprintf("repeated mismatch for field %s: %v vs %v", f1.Name, f1.Repeated, f2.Repeated))
			}
			if f1.Type == bigquery.RecordFieldType {
				f1.Fields = MergeSchemas(&Schema{Fields: f1.Fields}, &Schema{Fields: f2.Fields}).Fields
			}
			f1.Required = f1.Required && f2.Required
		} else {
			clone := *f2
			clone.Required = false
			fieldMap[clone.Name] = &clone
		}
	}

	for _, f1 := range s1.Fields {
		if _, ok := fieldMap[f1.Name]; !ok {
			f1.Required = false
		}
	}
	
	merged := &Schema{}
	for _, f := range fieldMap {
		merged.Fields = append(merged.Fields, f)
	}

	sort.Slice(merged.Fields, func(i, j int) bool {
		return merged.Fields[i].Name < merged.Fields[j].Name
	})

	return merged
}


func mustMarshal(value interface{}) []byte {
	bytes, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return bytes
}

// InferSchemaConcurrently infers the schema from a channel of JSON strings concurrently.
func InferSchemaConcurrently(lines <-chan string, numWorkers int) *Schema {
	var wg sync.WaitGroup
	schemas := make(chan *Schema, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var localMergedSchema *Schema
			for line := range lines {
				schema, err := InferSchema([]byte(line))
				if err != nil {
					// In a real application, you'd want to handle this error better.
					continue
				}
				if localMergedSchema == nil {
					localMergedSchema = schema
				} else {
					localMergedSchema = MergeSchemas(localMergedSchema, schema)
				}
			}
			schemas <- localMergedSchema
		}()
	}

	go func() {
		wg.Wait()
		close(schemas)
	}()

	var finalSchema *Schema
	for schema := range schemas {
		if schema == nil {
			continue
		}
		if finalSchema == nil {
			finalSchema = schema
		} else {
			finalSchema = MergeSchemas(finalSchema, schema)
		}
	}

	sort.Slice(finalSchema.Fields, func(i, j int) bool {
		return finalSchema.Fields[i].Name < finalSchema.Fields[j].Name
	})

	return finalSchema
}
