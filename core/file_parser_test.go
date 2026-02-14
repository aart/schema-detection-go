package core

import (
	"os"
	"testing"
)

func TestReadLines(t *testing.T) {
	// Create a temporary file with some content
	tmpfile, err := os.CreateTemp("", "test.ndjson")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	content := []byte("{\"name\": \"test1\"}\n{\"name\": \"test2\"}")
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	lines, err := ReadLines(tmpfile.Name())
	if err != nil {
		t.Fatalf("ReadLines failed: %v", err)
	}

	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}

	if lines[0] != "{\"name\": \"test1\"}" {
		t.Errorf("Expected line 1 to be '{\"name\": \"test1\"}', got '%s'", lines[0])
	}

	if lines[1] != "{\"name\": \"test2\"}" {
		t.Errorf("Expected line 2 to be '{\"name\": \"test2\"}', got '%s'", lines[1])
	}
}

func TestReadLines_NonExistentFile(t *testing.T) {
	_, err := ReadLines("nonexistent.ndjson")
	if err == nil {
		t.Fatal("Expected an error for a nonexistent file, but got nil")
	}
}
