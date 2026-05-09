package pipeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/baditaflorin/civitas/internal/evidence"
	"github.com/baditaflorin/civitas/internal/observability"
)

type fixtureExpectation struct {
	Shape                 string   `json:"shape"`
	State                 string   `json:"state"`
	MinConfidence         float64  `json:"min_confidence"`
	MinRowCount           int      `json:"min_row_count"`
	MinSizeBytes          int      `json:"min_size_bytes"`
	MaxEntities           int      `json:"max_entities"`
	RequiredFields        []string `json:"required_fields"`
	RequiredTimelineDates []string `json:"required_timeline_dates"`
	RequiredAnomalyCodes  []string `json:"required_anomaly_codes"`
	ForbiddenPrefix       string   `json:"forbidden_preview_prefix"`
}

func TestRealDataFixtures(t *testing.T) {
	dir := filepath.Join("..", "..", "test", "fixtures", "realdata")
	expectedFiles, err := filepath.Glob(filepath.Join(dir, "*.expected.json"))
	if err != nil {
		t.Fatalf("glob expectations: %v", err)
	}
	if len(expectedFiles) != 10 {
		t.Fatalf("expected 10 real-data expectations, got %d", len(expectedFiles))
	}
	pipe := New(DefaultRegistry(), observability.NewMetrics())
	for _, expectedPath := range expectedFiles {
		name := strings.TrimSuffix(filepath.Base(expectedPath), ".expected.json")
		t.Run(name, func(t *testing.T) {
			var expected fixtureExpectation
			body, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatalf("read expected: %v", err)
			}
			if err := json.Unmarshal(body, &expected); err != nil {
				t.Fatalf("decode expected: %v", err)
			}
			content, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				t.Fatalf("read fixture: %v", err)
			}
			doc := pipe.Analyze("case_realdata", "doc_"+stableID("fixture", name), name, content)
			assertFixture(t, doc, expected)
		})
	}
}

func TestRealDataFixturesAreDeterministic(t *testing.T) {
	dir := filepath.Join("..", "..", "test", "fixtures", "realdata")
	expectedFiles, err := filepath.Glob(filepath.Join(dir, "*.expected.json"))
	if err != nil {
		t.Fatalf("glob expectations: %v", err)
	}
	pipe := New(DefaultRegistry(), observability.NewMetrics())
	for _, expectedPath := range expectedFiles {
		name := strings.TrimSuffix(filepath.Base(expectedPath), ".expected.json")
		t.Run(name, func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				t.Fatalf("read fixture: %v", err)
			}
			docID := "doc_" + stableID("fixture", name)
			first := normalizedDeterminismDocument(pipe.Analyze("case_realdata", docID, name, content))
			second := normalizedDeterminismDocument(pipe.Analyze("case_realdata", docID, name, content))
			firstBody, err := json.Marshal(first)
			if err != nil {
				t.Fatalf("marshal first document: %v", err)
			}
			secondBody, err := json.Marshal(second)
			if err != nil {
				t.Fatalf("marshal second document: %v", err)
			}
			if string(firstBody) != string(secondBody) {
				t.Fatalf("analysis is not deterministic\nfirst:  %s\nsecond: %s", firstBody, secondBody)
			}
		})
	}
}

func assertFixture(t *testing.T, doc evidence.Document, expected fixtureExpectation) {
	t.Helper()
	if doc.Shape != expected.Shape {
		t.Fatalf("shape = %s, want %s", doc.Shape, expected.Shape)
	}
	if doc.State != expected.State {
		t.Fatalf("state = %s, want %s", doc.State, expected.State)
	}
	if doc.Confidence < expected.MinConfidence {
		t.Fatalf("confidence = %.2f, want >= %.2f", doc.Confidence, expected.MinConfidence)
	}
	if expected.MinRowCount > 0 && doc.Parse.RowCount < expected.MinRowCount {
		t.Fatalf("row count = %d, want >= %d", doc.Parse.RowCount, expected.MinRowCount)
	}
	if expected.MinSizeBytes > 0 && doc.Size < int64(expected.MinSizeBytes) {
		t.Fatalf("size = %d, want >= %d", doc.Size, expected.MinSizeBytes)
	}
	if (expected.MaxEntities > 0 || expected.Shape == "empty") && len(doc.Entities) > expected.MaxEntities {
		t.Fatalf("entities = %d, want <= %d", len(doc.Entities), expected.MaxEntities)
	}
	fieldNames := map[string]bool{}
	for _, field := range doc.Fields {
		fieldNames[field.Name] = true
	}
	for _, required := range expected.RequiredFields {
		if !fieldNames[required] {
			t.Fatalf("missing required field %q in %#v", required, doc.Fields)
		}
	}
	anomalyCodes := map[string]bool{}
	for _, anomaly := range doc.Anomalies {
		anomalyCodes[anomaly.Code] = true
		if anomaly.Message == "" || anomaly.Why == "" || anomaly.NextStep == "" {
			t.Fatalf("anomaly lacks actionable text: %#v", anomaly)
		}
	}
	for _, required := range expected.RequiredAnomalyCodes {
		if !anomalyCodes[required] {
			t.Fatalf("missing required anomaly %q in %#v", required, doc.Anomalies)
		}
	}
	dates := map[string]bool{}
	for _, event := range doc.Timeline {
		dates[event.When.Format("2006-01-02")] = true
	}
	for _, required := range expected.RequiredTimelineDates {
		if !dates[required] {
			t.Fatalf("missing timeline date %q in %#v", required, doc.Timeline)
		}
	}
	if expected.ForbiddenPrefix != "" && strings.HasPrefix(strings.TrimSpace(doc.Preview), expected.ForbiddenPrefix) {
		t.Fatalf("preview starts with forbidden prefix %q", expected.ForbiddenPrefix)
	}
}

func normalizedDeterminismDocument(doc evidence.Document) evidence.Document {
	doc.CreatedAt = time.Time{}
	doc.Parse.DurationMS = 0
	for i := range doc.Processors {
		doc.Processors[i].StartedAt = time.Time{}
		doc.Processors[i].EndedAt = time.Time{}
	}
	return doc
}
