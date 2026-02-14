package core

import (
	"os"
	"testing"
)

func TestReadLines(t *testing.T) {
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
	linesChan := ReadLines(paths, 2)

	var lines []string
	for line := range linesChan {
		lines = append(lines, line)
	}

	if len(lines) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(lines))
	}
}
