package main

import (
	"os"
	"testing"
)

func TestWriteJSONLine(t *testing.T) {
	f, err := os.CreateTemp("", "output*.jsonl")
	if err != nil {
		t.Fatalf("Could not create temp file: %v", err)
	}
	defer os.Remove(f.Name())

	data := PageData{
		URL:   "https://example.com",
		Title: "Example",
		Text:  "This is an example.",
		Tags:  []string{"example", "test"},
	}

	writeJSONLine(f, data)
	fi, err := f.Stat()
	if err != nil {
		t.Fatalf("Could not stat file: %v", err)
	}

	if fi.Size() == 0 {
		t.Error("Expected non-empty file")
	}
}

func TestSaveHTMLFile(t *testing.T) {
	dir := t.TempDir()
	content := []byte("<html><body>Example</body></html>")
	filename := "example.html"

	saveHTMLFile(dir, filename, content)
	filePath := dir + "/" + filename
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("Expected file to exist: %v", err)
	}

	savedContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Could not read saved file: %v", err)
	}

	if string(savedContent) != string(content) {
		t.Error("Saved content does not match expected content")
	}
}
