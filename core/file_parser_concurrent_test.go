package core

import (
	"os"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/google/go-cmp/cmp"
)

func TestReadLinesConcurrently(t *testing.T) {
	// Create two temporary files with some content
	tmpfile1, err := os.CreateTemp("", "test1.ndjson")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile1.Name())

	content1 := []byte("{\"name\": \"test1\"}\n{\"name\": \"test2\"}")
	if _, err := tmpfile1.Write(content1); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile1.Close(); err != nil {
		t.Fatal(err)
	}

	tmpfile2, err := os.CreateTemp("", "test2.ndjson")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile2.Name())

	content2 := []byte("{\"name\": \"test3\"}\n{\"name\": \"test4\"}")
	if _, err := tmpfile2.Write(content2); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile2.Close(); err != nil {
		t.Fatal(err)
	}

	paths := []string{tmpfile1.Name(), tmpfile2.Name()}
	linesChan := ReadLinesConcurrently(paths)

	var lines []string
	for line := range linesChan {
		lines = append(lines, line)
	}

	if len(lines) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(lines))
	}
}

func TestInferSchemaConcurrently(t *testing.T) {
	// Create two temporary files with some content
	tmpfile1, err := os.CreateTemp("", "test1.ndjson")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile1.Name())

	content1 := []byte("{\"name\": \"test1\", \"age\": 20}\n{\"name\": \"test2\", \"is_developer\": true}")
	if _, err := tmpfile1.Write(content1); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile1.Close(); err != nil {
		t.Fatal(err)
	}

	tmpfile2, err := os.CreateTemp("", "test2.ndjson")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile2.Name())

	content2 := []byte("{\"name\": \"test3\", \"age\": 30}\n{\"name\": \"test4\", \"city\": \"New York\"}")
	if _, err := tmpfile2.Write(content2); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile2.Close(); err != nil {
		t.Fatal(err)
	}

	paths := []string{tmpfile1.Name(), tmpfile2.Name()}
	linesChan := ReadLinesConcurrently(paths)
	finalSchema := InferSchemaConcurrently(linesChan, 2)

	expectedSchema := &Schema{
		Fields: []*FieldSchema{
			{Name: "age", Type: bigquery.IntegerFieldType, Required: false},
			{Name: "city", Type: bigquery.StringFieldType, Required: false},
			{Name: "is_developer", Type: bigquery.BooleanFieldType, Required: false},
			{Name: "name", Type: bigquery.StringFieldType, Required: true},
		},
	}

	if !cmp.Equal(expectedSchema, finalSchema) {
		t.Errorf("Expected schema %v, got %v", expectedSchema, finalSchema)
	}
}
