package pipeline

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/baditaflorin/civitas/internal/evidence"
	"github.com/baditaflorin/civitas/internal/observability"
)

const schemaVersion = "phase2.v1"

type Pipeline struct {
	registry Registry
	metrics  *observability.Metrics
}

type analysisResult struct {
	shape      string
	state      string
	text       string
	preview    string
	summary    string
	confidence float64
	fields     []evidence.FieldInference
	anomalies  []evidence.Anomaly
	timeline   []evidence.TimelineEvent
	parse      evidence.ParseMetrics
}

func New(registry Registry, metrics *observability.Metrics) *Pipeline {
	return &Pipeline{registry: registry, metrics: metrics}
}

func (p *Pipeline) Tools() []Tool {
	return p.registry.Tools()
}

func (p *Pipeline) Analyze(caseID, docID, filename string, content []byte) evidence.Document {
	started := time.Now().UTC()
	p.metrics.IngestionJobsStarted.Inc()
	contentType := detectContentType(content)
	hash := sha256.Sum256(content)
	hashText := hex.EncodeToString(hash[:])

	result := p.analyzeContent(docID, filename, contentType, content)
	entities := extractEntities(docID, entityExtractionText(result))
	processors := p.processorStatuses(started, result.shape, result.state)

	doc := evidence.Document{
		ID:          docID,
		CaseID:      caseID,
		Filename:    filepath.Base(filename),
		ContentType: contentType,
		Size:        int64(len(content)),
		SHA256:      hashText,
		Status:      result.state,
		Shape:       result.shape,
		State:       result.state,
		Confidence:  result.confidence,
		Text:        result.text,
		Preview:     result.preview,
		Summary:     result.summary,
		Fields:      sortedFields(result.fields),
		Anomalies:   sortedAnomalies(result.anomalies),
		Entities:    entities,
		Timeline:    result.timeline,
		Processors:  processors,
		Provenance: evidence.Provenance{
			SchemaVersion: schemaVersion,
			AppVersion:    "0.3.0",
			SourceID:      docID,
			SourceName:    filepath.Base(filename),
			SourceSHA256:  hashText,
			Parameters:    []string{"phase2-classifier", "deterministic-heuristics"},
		},
		Parse:     result.parse,
		CreatedAt: started,
	}
	p.metrics.DocumentsProcessed.Inc()
	p.metrics.IngestionJobsCompleted.Inc()
	return doc
}

func detectContentType(content []byte) string {
	if len(content) == 0 {
		return "application/x-empty"
	}
	return http.DetectContentType(content)
}

func (p *Pipeline) analyzeContent(docID, filename, contentType string, content []byte) analysisResult {
	begin := time.Now()
	shape := classifyShape(filename, contentType, content)
	result := analyzeByShape(docID, filename, contentType, shape, content)
	result.shape = shape
	result.parse.DurationMS = time.Since(begin).Milliseconds()
	result.parse.SizeBucket = sizeBucket(len(content))
	if result.preview == "" {
		result.preview = preview(result.text)
	}
	if result.summary == "" {
		result.summary = summarize(filename, result)
	}
	result.timeline = sortedTimeline(append(result.timeline, extractTimeline(docID, result.text)...))
	result.anomalies = append(result.anomalies, processorAnomalies(shape, p.registry)...)
	if result.state != "failed" {
		result.confidence = clampConfidence(result.confidence, result.anomalies)
	}
	return result
}

func (p *Pipeline) processorStatuses(started time.Time, shape, state string) []evidence.ProcessorStatus {
	tools := p.registry.Tools()
	statuses := make([]evidence.ProcessorStatus, 0, len(tools)+1)
	for _, tool := range tools {
		status := "available"
		msg := "native command discovered"
		if !tool.Available {
			status = "missing"
			msg = "native command not installed"
		}
		statuses = append(statuses, evidence.ProcessorStatus{
			Name:      tool.Name,
			Kind:      tool.Kind,
			Available: tool.Available,
			Status:    status,
			Message:   msg,
			StartedAt: started,
			EndedAt:   time.Now().UTC(),
		})
	}
	statuses = append(statuses, evidence.ProcessorStatus{
		Name:      "Civitas Phase 2 inference engine",
		Kind:      "classification-normalization-confidence",
		Available: true,
		Status:    state,
		Message:   "classified evidence as " + shape,
		StartedAt: started,
		EndedAt:   time.Now().UTC(),
	})
	return statuses
}

func summarize(filename string, result analysisResult) string {
	parts := []string{filepath.Base(filename), result.shape, result.state}
	if result.preview != "" {
		parts = append(parts, preview(result.preview))
	}
	if len(result.anomalies) > 0 {
		parts = append(parts, result.anomalies[0].Message)
	}
	return strings.Join(parts, ": ")
}

func entityExtractionText(result analysisResult) string {
	if result.shape == "csv" {
		return result.preview + " " + strings.Join(fieldValues(result.fields), " ")
	}
	if len(result.text) > 200_000 {
		return previewLong(result.text, 200_000)
	}
	return result.text
}
