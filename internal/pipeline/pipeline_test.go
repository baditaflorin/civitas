package pipeline

import (
	"strings"
	"testing"

	"github.com/baditaflorin/civitas/internal/observability"
)

func TestAnalyzeExtractsInvestigationSignals(t *testing.T) {
	pipe := New(DefaultRegistry(), observability.NewMetrics())
	doc := pipe.Analyze(
		"case_test",
		"doc_test",
		"lead.txt",
		[]byte("Payment due on 2026-05-08 from source@example.org at 42 Civic Street for EUR 1200."),
	)

	if doc.State != "ready" {
		t.Fatalf("expected ready document, got %s", doc.State)
	}
	if !strings.Contains(doc.Text, "source@example.org") {
		t.Fatalf("expected extracted text to contain email, got %q", doc.Text)
	}
	if len(doc.Entities) < 3 {
		t.Fatalf("expected email, address, and money entities, got %#v", doc.Entities)
	}
	if len(doc.Timeline) != 1 {
		t.Fatalf("expected one timeline event, got %d", len(doc.Timeline))
	}
	if len(doc.Processors) == 0 {
		t.Fatal("expected processor statuses")
	}
}
