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
		fieldType, err := inferType(value)
		if err != nil {
			return nil, err
		}
		schema.Fields = append(schema.Fields, &FieldSchema{
			Name:     key,
			Type:     fieldType,
			Required: true,
		})
	}

	sort.Slice(schema.Fields, func(i, j int) bool {
		return schema.Fields[i].Name < schema.Fields[j].Name
	})

	return schema, nil
}

func inferType(value interface{}) (bigquery.FieldType, error) {
	switch value.(type) {
	case string:
		return bigquery.StringFieldType, nil
	case float64:
		// json.Unmarshal uses float64 for all numbers.
		// Check if the number is an integer.
		if float64(int64(value.(float64))) == value.(float64) {
			return bigquery.IntegerFieldType, nil
		}
		return bigquery.FloatFieldType, nil
	case bool:
		return bigquery.BooleanFieldType, nil
	default:
		return "", fmt.Errorf("unsupported type: %T", value)
	}
}
