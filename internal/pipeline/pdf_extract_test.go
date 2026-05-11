package pipeline

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractPDFText_RealDocument(t *testing.T) {
	path := filepath.Join("..", "..", "test", "fixtures", "realdata", "blm_foia_sample.pdf")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	text, pages, err := extractPDFText(content)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	if pages != 2 {
		t.Fatalf("pages = %d, want 2", pages)
	}
	if !strings.Contains(text, "Freedom of Information Act") {
		t.Fatalf("expected FOIA phrase in extracted text, got: %q", text)
	}
}

func TestExtractPDFText_GarbageReturnsError(t *testing.T) {
	// %PDF header but nothing else valid; the parser should fail cleanly,
	// not panic, and the wrapper should turn either into a returned error.
	_, _, err := extractPDFText([]byte("%PDF-1.4\nnot a real pdf\n%%EOF"))
	if err == nil {
		t.Fatal("expected error from garbage PDF, got nil")
	}
	if errors.Is(err, errPDFNoText) {
		t.Fatalf("expected parse error, got no-text error")
	}
}
