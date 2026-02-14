package core

import (
	"encoding/json"
	"fmt"
	"sort"

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
		schema.Fields = append(schema.Fields, fieldSchema)
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
					mergedSchema = mergeSchemas(mergedSchema, itemSchema)
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
	default:
		return nil, fmt.Errorf("unsupported type: %T", value)
	}
}

func mergeSchemas(s1, s2 *Schema) *Schema {
	merged := &Schema{Fields: s1.Fields}
	for _, f2 := range s2.Fields {
		found := false
		for _, f1 := range merged.Fields {
			if f1.Name == f2.Name {
				found = true
				if f1.Type != f2.Type {
					// In a real scenario, you might want to handle this more gracefully,
					// e.g., by picking a more general type.
					// For now, we'll just panic.
					panic(fmt.Sprintf("type mismatch for field %s: %s vs %s", f1.Name, f1.Type, f2.Type))
				}
				if f1.Repeated != f2.Repeated {
					panic(fmt.Sprintf("repeated mismatch for field %s: %v vs %v", f1.Name, f1.Repeated, f2.Repeated))
				}
				if f1.Type == bigquery.RecordFieldType {
					f1.Fields = mergeSchemas(&Schema{Fields: f1.Fields}, &Schema{Fields: f2.Fields}).Fields
				}
			}
		}
		if !found {
			f2.Required = false
			merged.Fields = append(merged.Fields, f2)
		}
	}

	// After merging, check if fields from s1 are present in s2. If not, mark them as not required.
	for _, f1 := range merged.Fields {
		found := false
		for _, f2 := range s2.Fields {
			if f1.Name == f2.Name {
				found = true
				break
			}
		}
		if !found {
			f1.Required = false
		}
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

