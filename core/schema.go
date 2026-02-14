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

func mustMarshal(value interface{}) []byte {
	bytes, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return bytes
}

