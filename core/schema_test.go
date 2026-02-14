package core

import (
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/google/go-cmp/cmp"
)

func TestInferSchema(t *testing.T) {
	testCases := []struct {
		name           string
		jsonData       string
		expectedSchema *Schema
	}{
		{
			name:     "simple case",
			jsonData: `{"name": "test", "age": 42, "is_developer": true}`,
			expectedSchema: &Schema{
				Fields: []*FieldSchema{
					{Name: "age", Type: bigquery.IntegerFieldType, Required: true},
					{Name: "is_developer", Type: bigquery.BooleanFieldType, Required: true},
					{Name: "name", Type: bigquery.StringFieldType, Required: true},
				},
			},
		},
		{
			name:     "float number",
			jsonData: `{"price": 12.34}`,
			expectedSchema: &Schema{
				Fields: []*FieldSchema{
					{Name: "price", Type: bigquery.FloatFieldType, Required: true},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inferredSchema, err := InferSchema([]byte(tc.jsonData))
			if err != nil {
				t.Fatalf("InferSchema failed: %v", err)
			}

			if !cmp.Equal(tc.expectedSchema, inferredSchema) {
				t.Errorf("Expected schema %v, got %v", tc.expectedSchema, inferredSchema)
			}
		})
	}
}

func TestInferSchema_UnsupportedType(t *testing.T) {
	jsonData := `{"data": null}`
	_, err := InferSchema([]byte(jsonData))
	if err == nil {
		t.Fatal("Expected an error for an unsupported type, but got nil")
	}
}
