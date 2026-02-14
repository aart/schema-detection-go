package core

import (
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/google/go-cmp/cmp"
)

func TestInferSchema(t *testing.T) {
	jsonData := `{"name": "test", "age": 42, "is_developer": true}`

	expectedSchema := &Schema{
		Fields: []*FieldSchema{
			{Name: "age", Type: bigquery.IntegerFieldType, Required: true},
			{Name: "is_developer", Type: bigquery.BooleanFieldType, Required: true},
			{Name: "name", Type: bigquery.StringFieldType, Required: true},
		},
	}

	inferredSchema, err := InferSchema([]byte(jsonData))
	if err != nil {
		t.Fatalf("InferSchema failed: %v", err)
	}

	if !cmp.Equal(expectedSchema, inferredSchema) {
		t.Errorf("Expected schema %v, got %v", expectedSchema, inferredSchema)
	}
}
