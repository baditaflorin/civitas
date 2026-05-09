package storage

import (
	"strings"
	"testing"
	"time"

	"github.com/baditaflorin/civitas/internal/evidence"
	"github.com/baditaflorin/civitas/internal/exporter"
)

func TestStoreCaseDocumentSearchGraphTimelineAndExport(t *testing.T) {
	store, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("new store: %v", err)
	}

	item, err := store.CreateCase("Harbor contracts", "public procurement")
	if err != nil {
		t.Fatalf("create case: %v", err)
	}

	doc := evidence.Document{
		ID:          "doc_test",
		CaseID:      item.ID,
		Filename:    "contract.txt",
		ContentType: "text/plain",
		Size:        12,
		SHA256:      "abc",
		Status:      "completed",
		Text:        "source@example.org signed the contract on 2026-05-08.",
		Entities: []evidence.Entity{{
			ID:         "email_source",
			Type:       "email",
			Value:      "source@example.org",
			Confidence: 0.9,
			Source:     "doc_test",
		}},
		Timeline: []evidence.TimelineEvent{{
			ID:         "event_test",
			Document:   "doc_test",
			When:       time.Date(2026, 5, 8, 0, 0, 0, 0, time.UTC),
			Label:      "Signed",
			Confidence: 0.8,
		}},
		CreatedAt: time.Now().UTC(),
	}

	if err := store.AddDocument(item.ID, doc, []byte(doc.Text)); err != nil {
		t.Fatalf("add document: %v", err)
	}

	results, err := store.Search(item.ID, "signed")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(results) != 1 || results[0].DocumentID != doc.ID {
		t.Fatalf("unexpected search results: %#v", results)
	}

	graph, err := store.Graph(item.ID)
	if err != nil {
		t.Fatalf("graph: %v", err)
	}
	if len(graph.Nodes) != 2 || len(graph.Edges) != 1 {
		t.Fatalf("unexpected graph: %#v", graph)
	}

	events, err := store.Timeline(item.ID)
	if err != nil {
		t.Fatalf("timeline: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("unexpected timeline: %#v", events)
	}

	savedExport, err := store.SaveExport(evidence.Export{
		ID:        "export_test",
		CaseID:    item.ID,
		Format:    "markdown",
		Body:      "# Export",
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("save export: %v", err)
	}
	if savedExport.Path == "" {
		t.Fatal("saved export path is empty")
	}

	exported, err := store.ReadExport(item.ID, "export_test")
	if err != nil {
		t.Fatalf("read export: %v", err)
	}
	if !strings.Contains(exported.Body, "Export") {
		t.Fatalf("unexpected export body: %q", exported.Body)
	}

	content, err := store.DocumentContent(item.ID, doc.ID)
	if err != nil {
		t.Fatalf("document content: %v", err)
	}
	refreshedCase, err := store.GetCase(item.ID)
	if err != nil {
		t.Fatalf("get refreshed case: %v", err)
	}
	state := exporter.CaseState(refreshedCase, []evidence.Document{doc}, map[string][]byte{doc.ID: content}, "0.3.0", time.Now().UTC())
	restoredStore, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("new restored store: %v", err)
	}
	imported, err := restoredStore.ImportCaseState(state)
	if err != nil {
		t.Fatalf("import case state: %v", err)
	}
	importedDocs, err := restoredStore.Documents(imported.ID)
	if err != nil {
		t.Fatalf("imported documents: %v", err)
	}
	if len(importedDocs) != 1 || importedDocs[0].ID != doc.ID {
		t.Fatalf("unexpected imported docs: %#v", importedDocs)
	}
	importedContent, err := restoredStore.DocumentContent(imported.ID, doc.ID)
	if err != nil {
		t.Fatalf("imported content: %v", err)
	}
	if string(importedContent) != doc.Text {
		t.Fatalf("unexpected imported content: %q", importedContent)
	}
}
