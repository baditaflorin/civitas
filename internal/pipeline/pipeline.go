package pipeline

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/baditaflorin/civitas/internal/evidence"
	"github.com/baditaflorin/civitas/internal/observability"
)

type Pipeline struct {
	registry Registry
	metrics  *observability.Metrics
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
	contentType := http.DetectContentType(content)
	hash := sha256.Sum256(content)

	text := extractReadableText(content, contentType)
	entities := extractEntities(docID, text)
	timeline := extractTimeline(docID, text)
	processors := p.processorStatuses(started)

	doc := evidence.Document{
		ID:          docID,
		CaseID:      caseID,
		Filename:    filename,
		ContentType: contentType,
		Size:        int64(len(content)),
		SHA256:      hex.EncodeToString(hash[:]),
		Status:      "completed",
		Text:        text,
		Summary:     summarize(filename, text, entities),
		Entities:    entities,
		Timeline:    timeline,
		Processors:  processors,
		CreatedAt:   started,
	}
	p.metrics.DocumentsProcessed.Inc()
	p.metrics.IngestionJobsCompleted.Inc()
	return doc
}

func (p *Pipeline) processorStatuses(started time.Time) []evidence.ProcessorStatus {
	tools := p.registry.Tools()
	statuses := make([]evidence.ProcessorStatus, 0, len(tools)+1)
	for _, tool := range tools {
		status := "skipped"
		msg := "adapter boundary present; native command not installed"
		if tool.Available {
			status = "available"
			msg = "native command discovered"
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
		Name:      "Civitas fallback extractor",
		Kind:      "text-entities-timeline",
		Available: true,
		Status:    "completed",
		Message:   "processed text-compatible content and generated v1 investigative signals",
		StartedAt: started,
		EndedAt:   time.Now().UTC(),
	})
	return statuses
}

func extractReadableText(content []byte, contentType string) string {
	if len(content) == 0 {
		return ""
	}
	text := string(content)
	if utf8.Valid(content) && printableRatio(text) > 0.75 {
		return strings.TrimSpace(text)
	}
	if strings.Contains(contentType, "pdf") {
		return "PDF evidence uploaded. Configure Tika, PyMuPDF, pdfminer, qpdf, or Ghostscript in the backend image for full text extraction."
	}
	if strings.HasPrefix(contentType, "image/") {
		return "Image evidence uploaded. Configure Tesseract, ImageMagick, dlib, or MediaPipe in the backend image for OCR and redaction."
	}
	if strings.HasPrefix(contentType, "audio/") || strings.HasPrefix(contentType, "video/") {
		return "Media evidence uploaded. Configure Whisper.cpp, pyannote, and NLLB-200 in the backend image for transcription and translation."
	}
	return fmt.Sprintf("Binary evidence uploaded with detected content type %s.", contentType)
}

func printableRatio(text string) float64 {
	if text == "" {
		return 0
	}
	printable := 0
	total := 0
	for _, r := range text {
		total++
		if r == '\n' || r == '\r' || r == '\t' || (r >= 32 && r < 127) {
			printable++
		}
	}
	return float64(printable) / float64(total)
}

func extractEntities(docID, text string) []evidence.Entity {
	patterns := map[string]*regexp.Regexp{
		"email":   regexp.MustCompile(`(?i)\b[A-Z0-9._%+\-]+@[A-Z0-9.\-]+\.[A-Z]{2,}\b`),
		"phone":   regexp.MustCompile(`\b(?:\+?\d[\d .()\-]{7,}\d)\b`),
		"money":   regexp.MustCompile(`(?i)\b(?:EUR|USD|RON|GBP|\$|€)\s?[0-9][0-9,.\s]*\b`),
		"address": regexp.MustCompile(`(?im)\b\d{1,5}\s+[A-Za-z0-9 .'\-]+(?:Street|St|Road|Rd|Avenue|Ave|Boulevard|Blvd|Lane|Ln|Square|Sq)\b`),
	}
	seen := map[string]bool{}
	var out []evidence.Entity
	for kind, pattern := range patterns {
		for _, match := range pattern.FindAllString(text, -1) {
			value := strings.TrimSpace(match)
			key := kind + ":" + strings.ToLower(value)
			if seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, evidence.Entity{
				ID:         stableID(kind, value),
				Type:       kind,
				Value:      value,
				Confidence: 0.72,
				Source:     docID,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Type == out[j].Type {
			return out[i].Value < out[j].Value
		}
		return out[i].Type < out[j].Type
	})
	return out
}

func extractTimeline(docID, text string) []evidence.TimelineEvent {
	datePattern := regexp.MustCompile(`\b(20\d{2}|19\d{2})-(0[1-9]|1[0-2])-([0-2]\d|3[01])\b`)
	matches := datePattern.FindAllString(text, -1)
	seen := map[string]bool{}
	var events []evidence.TimelineEvent
	for _, match := range matches {
		if seen[match] {
			continue
		}
		seen[match] = true
		when, err := time.Parse("2006-01-02", match)
		if err != nil {
			continue
		}
		events = append(events, evidence.TimelineEvent{
			ID:         stableID("event", docID+match),
			Document:   docID,
			When:       when,
			Label:      "Mentioned date " + match,
			Confidence: 0.7,
		})
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].When.Before(events[j].When)
	})
	return events
}

func summarize(filename, text string, entities []evidence.Entity) string {
	if text == "" {
		return filename + " contains no extracted text yet."
	}
	clean := strings.Join(strings.Fields(text), " ")
	if len(clean) > 220 {
		clean = clean[:220] + "..."
	}
	return fmt.Sprintf("%s: %s (%d entities)", filename, clean, len(entities))
}

func stableID(prefix, value string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(value)))
	return prefix + "_" + hex.EncodeToString(sum[:])[:16]
}
