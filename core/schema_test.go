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
		{
			name:     "nested record",
			jsonData: `{"person": {"name": "test", "age": 30}}`,
			expectedSchema: &Schema{
				Fields: []*FieldSchema{
					{
						Name: "person",
						Type: bigquery.RecordFieldType,
						Fields: []*FieldSchema{
							{Name: "age", Type: bigquery.IntegerFieldType, Required: true},
							{Name: "name", Type: bigquery.StringFieldType, Required: true},
						},
						Required: true,
					},
				},
			},
		},
		{
			name:     "repeated primitive",
			jsonData: `{"scores": [1, 2, 3]}`,
			expectedSchema: &Schema{
				Fields: []*FieldSchema{
					{Name: "scores", Type: bigquery.IntegerFieldType, Repeated: true, Required: true},
				},
			},
		},
		{
			name:     "repeated record",
			jsonData: `{"people": [{"name": "test1"}, {"name": "test2", "age": 30}]}`,
			expectedSchema: &Schema{
				Fields: []*FieldSchema{
					{
						Name: "people",
						Type: bigquery.RecordFieldType,
						Repeated: true,
						Fields: []*FieldSchema{
							{Name: "age", Type: bigquery.IntegerFieldType, Required: false},
							{Name: "name", Type: bigquery.StringFieldType, Required: true},
						},
						Required: true,
					},
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

func TestInferSchema_NullableField(t *testing.T) {
	jsonLines := []string{
		`{"name": "test1"}`,
		`{"name": "test2", "age": 30}`,
	}

	var mergedSchema *Schema
	for _, line := range jsonLines {
		schema, err := InferSchema([]byte(line))
		if err != nil {
			t.Fatalf("InferSchema failed: %v", err)
		}
		if mergedSchema == nil {
			mergedSchema = schema
		} else {
			mergedSchema = MergeSchemas(mergedSchema, schema)
		}
	}

	expectedSchema := &Schema{
		Fields: []*FieldSchema{
			{Name: "age", Type: bigquery.IntegerFieldType, Required: false},
			{Name: "name", Type: bigquery.StringFieldType, Required: true},
		},
	}

	if !cmp.Equal(expectedSchema, mergedSchema) {
		t.Errorf("Expected schema %v, got %v", expectedSchema, mergedSchema)
	}
}

func TestInferSchema_NullValue(t *testing.T) {
	jsonLines := []string{
		`{"name": "test1", "age": 30}`,
		`{"name": "test2", "age": null}`,
	}

	var mergedSchema *Schema
	for _, line := range jsonLines {
		schema, err := InferSchema([]byte(line))
		if err != nil {
			t.Fatalf("InferSchema failed: %v", err)
		}
		if mergedSchema == nil {
			mergedSchema = schema
		} else {
			mergedSchema = MergeSchemas(mergedSchema, schema)
		}
	}

	expectedSchema := &Schema{
		Fields: []*FieldSchema{
			{Name: "age", Type: bigquery.IntegerFieldType, Required: true},
			{Name: "name", Type: bigquery.StringFieldType, Required: true},
		},
	}

	if !cmp.Equal(expectedSchema, mergedSchema) {
		t.Errorf("Expected schema %v, got %v", expectedSchema, mergedSchema)
	}
}

func TestInferSchema_MalformedJSON(t *testing.T) {
	testCases := []string{
		`{"name": "test", "age": 42,}`,
		`{"name": "test", "age": 42`,
	}

	for _, tc := range testCases {
		_, err := InferSchema([]byte(tc))
		if err == nil {
			t.Errorf("Expected an error for malformed JSON '%s', but got nil", tc)
		}
	}
}

func TestInferSchema_UnsupportedType(t *testing.T) {
	jsonData := `{"data": null}`
	schema, err := InferSchema([]byte(jsonData))
	if err != nil {
		t.Fatalf("InferSchema failed: %v", err)
	}
	if len(schema.Fields) > 0 {
		t.Errorf("Expected 0 fields, got %d", len(schema.Fields))
	}
}
