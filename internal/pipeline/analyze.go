package pipeline

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/baditaflorin/civitas/internal/evidence"
	"golang.org/x/net/html"
)

func analyzeByShape(docID, filename, contentType, shape string, content []byte) analysisResult {
	switch shape {
	case "empty":
		return emptyResult(docID)
	case "csv":
		return analyzeCSV(docID, filename, content)
	case "ocds_json", "json":
		return analyzeJSON(docID, filename, content, shape)
	case "html_article", "html_data_source":
		return analyzeHTML(docID, filename, content, shape)
	case "text":
		text := normalizeBytes(content)
		return analysisResult{state: "ready", text: text, preview: preview(text), confidence: 0.8}
	case "pdf":
		return analyzePDF(docID, filename, content)
	case "image_scan":
		return processorNeeded(docID, filename, shape, "ocr_processor_unavailable", "OCR is required for this scan.", "Install or enable Tesseract/ImageMagick/MediaPipe, then reprocess this evidence.", 0.78)
	case "audio":
		return processorNeeded(docID, filename, shape, "transcription_processor_unavailable", "Transcription is required for this audio.", "Install or enable Whisper.cpp/pyannote, then reprocess this evidence.", 0.78)
	case "archive_zip":
		return analyzeArchive(docID, filename, content)
	default:
		text := normalizeText(string(content))
		return analysisResult{
			state:      "unsupported",
			text:       text,
			preview:    preview(text),
			confidence: 0.45,
			anomalies:  []evidence.Anomaly{newAnomaly(docID, "unsupported_binary", "warning", "Civitas stored this file but cannot infer its evidence shape.", "The content type is "+contentType+" and no v2 parser matched it.", "Keep it as source evidence or convert it to PDF, CSV, JSON, HTML, image, audio, or ZIP.", 0.8)},
		}
	}
}

func emptyResult(docID string) analysisResult {
	return analysisResult{
		state:      "failed",
		confidence: 1,
		anomalies: []evidence.Anomaly{
			newAnomaly(docID, "empty_file", "error", "This upload is empty.", "The file has zero bytes, which usually means a failed transfer or placeholder file.", "Re-upload the original file or keep this record as a failed transfer.", 1),
		},
	}
}

func analyzeCSV(docID, filename string, content []byte) analysisResult {
	_ = filename
	text := normalizeBytes(content)
	delimiter := sniffDelimiter(text)
	reader := csv.NewReader(strings.NewReader(text))
	reader.Comma = delimiter
	reader.FieldsPerRecord = -1
	reader.ReuseRecord = true
	rows := [][]string{}
	fieldCounts := map[int]int{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return analysisResult{
				state:      "recoverable_error",
				text:       text,
				preview:    preview(text),
				confidence: 0.55,
				anomalies:  []evidence.Anomaly{newAnomaly(docID, "csv_parse_error", "error", "CSV parsing stopped before the file was fully understood.", err.Error(), "Check quoting, embedded newlines, or export the spreadsheet as UTF-8 CSV.", 0.9)},
			}
		}
		cp := append([]string(nil), record...)
		if len(rows) < 200 {
			rows = append(rows, cp)
		}
		fieldCounts[len(record)]++
	}
	rowCount := totalCounts(fieldCounts)
	fieldCount := modalCount(fieldCounts)
	headers, hasHeader := inferHeaders(rows, fieldCount)
	fields := inferCSVFields(docID, headers, rows, hasHeader)
	anomalies := []evidence.Anomaly{}
	if !hasHeader {
		anomalies = append(anomalies, newAnomaly(docID, "csv_header_missing", "warning", "CSV rows were detected but no header row was found.", "The first row looks like data, so Civitas assigned stable column names.", "Rename columns in the export or provide a CSV with headers for stronger field names.", 0.86))
	}
	if len(fieldCounts) > 1 {
		anomalies = append(anomalies, newAnomaly(docID, "csv_mixed_field_counts", "warning", "Some CSV rows have different field counts.", "Rows do not all share the modal field count "+strconv.Itoa(fieldCount)+".", "Review rows with unusual delimiters, quotes, or embedded newlines.", 0.78))
	}
	return analysisResult{
		state:      "ready",
		text:       text,
		preview:    csvPreview(headers, rows, hasHeader),
		confidence: 0.82,
		fields:     fields,
		anomalies:  anomalies,
		parse: evidence.ParseMetrics{
			RowCount:   rowCount,
			FieldCount: fieldCount,
		},
	}
}

func analyzeJSON(docID, filename string, content []byte, shape string) analysisResult {
	_ = filename
	var obj any
	if err := json.Unmarshal(content, &obj); err != nil {
		return analysisResult{state: "recoverable_error", confidence: 0.65, text: normalizeBytes(content), anomalies: []evidence.Anomaly{newAnomaly(docID, "json_parse_error", "error", "JSON could not be parsed.", err.Error(), "Check for truncation, trailing commas, or invalid encoding.", 0.9)}}
	}
	pretty, _ := json.MarshalIndent(obj, "", "  ")
	text := string(pretty)
	fields := inferJSONFields(docID, obj, shape)
	return analysisResult{
		state:      "ready",
		text:       text,
		preview:    preview(strings.Join(fieldValues(fields), " | ")),
		confidence: 0.94,
		fields:     fields,
		timeline:   timelineFromFields(docID, fields),
	}
}

func analyzeHTML(docID, filename string, content []byte, shape string) analysisResult {
	_ = filename
	_ = shape
	text, title, times, links := extractHTML(content)
	fields := []evidence.FieldInference{}
	if title != "" {
		fields = append(fields, newField(docID, "title", "text", title, title, 0.88, "HTML title or heading"))
	}
	if len(times) > 0 {
		if normalized, ok := normalizeDate(times[0]); ok {
			fields = append(fields, newField(docID, "published_date", "date", times[0], normalized, 0.84, "HTML time element"))
		}
	}
	if len(links) > 0 {
		fields = append(fields, newField(docID, "source_url", "url", links[0], links[0], 0.7, "First source link on page"))
	}
	return analysisResult{
		state:      "ready",
		text:       text,
		preview:    preview(text),
		confidence: 0.82,
		fields:     fields,
		timeline:   timelineFromFields(docID, fields),
	}
}

func analyzePDF(docID, filename string, content []byte) analysisResult {
	if !bytes.HasPrefix(content, []byte("%PDF")) || !bytes.Contains(content, []byte("%%EOF")) || len(content) < 512 {
		return analysisResult{
			state:      "recoverable_error",
			text:       "PDF evidence appears partial or corrupt.",
			preview:    "PDF evidence appears partial or corrupt.",
			confidence: 0.86,
			anomalies:  []evidence.Anomaly{newAnomaly(docID, "pdf_truncated_or_corrupt", "error", "PDF appears truncated or corrupt.", "The file is missing a reliable PDF header/end marker or is too small for the expected document.", "Re-upload the original file or run qpdf repair before reprocessing.", 0.9)},
		}
	}
	text, pages, err := extractPDFText(content)
	switch {
	case err == nil:
		normalized := normalizeText(text)
		return analysisResult{
			state:      "ready",
			text:       normalized,
			preview:    preview(normalized),
			confidence: 0.82,
			fields: []evidence.FieldInference{
				{Name: "pages", Value: strconv.Itoa(pages), Confidence: 0.95},
			},
		}
	case errors.Is(err, errPDFNoText):
		// Parsed fine, but no text layer — this is a scanned PDF and needs OCR.
		return processorNeeded(docID, filename, "pdf", "ocr_processor_unavailable", "PDF has no text layer; OCR is required.", "Install or enable Tesseract or another OCR backend, then reprocess this PDF.", 0.7)
	default:
		// Parse error: the file is structurally a PDF but our parser cannot read it.
		return analysisResult{
			state:      "recoverable_error",
			text:       "PDF could not be parsed: " + err.Error(),
			preview:    "PDF could not be parsed.",
			confidence: 0.7,
			anomalies:  []evidence.Anomaly{newAnomaly(docID, "pdf_parse_failed", "error", "PDF parser rejected the file.", err.Error(), "Try repairing the file with qpdf, or re-export from the source application.", 0.85)},
		}
	}
}

func analyzeArchive(docID, filename string, content []byte) analysisResult {
	_ = filename
	if err := validateZip(content); err != nil {
		return analysisResult{
			state:      "recoverable_error",
			text:       "ZIP archive appears partial or corrupt.",
			preview:    "ZIP archive appears partial or corrupt.",
			confidence: 0.86,
			anomalies:  []evidence.Anomaly{newAnomaly(docID, "archive_truncated_or_corrupt", "error", "Archive appears truncated or corrupt.", err.Error(), "Re-download or re-upload the complete archive; partial recovery may be possible with archive tools.", 0.92)},
		}
	}
	return analysisResult{state: "ready", text: "ZIP archive validated.", preview: "ZIP archive validated.", confidence: 0.9}
}

func processorNeeded(docID, filename, shape, code, message, nextStep string, confidence float64) analysisResult {
	_ = filename
	text := message
	return analysisResult{
		state:      "needs_processor",
		text:       text,
		preview:    text,
		confidence: confidence,
		anomalies:  []evidence.Anomaly{newAnomaly(docID, code, "warning", message, "The file is valid "+shape+" evidence, but the required native processor is not available in this runtime.", nextStep, confidence)},
	}
}

func extractHTML(content []byte) (string, string, []string, []string) {
	root, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return normalizeBytes(content), "", nil, nil
	}
	var title string
	var texts, times, links []string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "script", "style", "nav", "footer", "header":
				return
			case "a":
				for _, attr := range n.Attr {
					if attr.Key == "href" && strings.HasPrefix(attr.Val, "http") {
						links = append(links, attr.Val)
					}
				}
			}
		}
		if n.Type == html.TextNode {
			value := strings.TrimSpace(n.Data)
			if value != "" {
				texts = append(texts, value)
				if title == "" && len(value) > 20 && len(value) < 180 {
					title = value
				}
				if _, ok := normalizeDate(value); ok {
					times = append(times, value)
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(root)
	return previewLong(strings.Join(texts, " "), 8000), title, uniqueStrings(times), uniqueStrings(links)
}

func newField(docID, name, kind, value, normalized string, confidence float64, reason string) evidence.FieldInference {
	return evidence.FieldInference{ID: stableID("field", docID+name+value), Name: name, Type: kind, Value: value, Normalized: normalized, Confidence: confidence, Reason: reason}
}

func newAnomaly(docID, code, severity, message, why, nextStep string, confidence float64) evidence.Anomaly {
	return evidence.Anomaly{ID: stableID("anomaly", docID+code+message), Code: code, Severity: severity, Message: message, Why: why, NextStep: nextStep, Confidence: confidence}
}

func totalCounts(counts map[int]int) int {
	total := 0
	for _, count := range counts {
		total += count
	}
	return total
}

func modalCount(counts map[int]int) int {
	best, bestN := 0, -1
	for fields, count := range counts {
		if count > bestN || (count == bestN && fields > best) {
			best, bestN = fields, count
		}
	}
	return best
}

func fieldValues(fields []evidence.FieldInference) []string {
	values := make([]string, 0, len(fields))
	for _, field := range fields {
		values = append(values, field.Name+"="+field.Value)
	}
	sort.Strings(values)
	return values
}

func processorAnomalies(shape string, registry Registry) []evidence.Anomaly {
	_ = registry
	_ = shape
	return nil
}

func clampConfidence(confidence float64, anomalies []evidence.Anomaly) float64 {
	for _, anomaly := range anomalies {
		if anomaly.Severity == "error" && confidence > 0.86 {
			confidence = 0.86
		}
	}
	return confidence
}

func sizeBucket(size int) string {
	switch {
	case size == 0:
		return "empty"
	case size < 1_000_000:
		return "small"
	case size < 5_000_000:
		return "medium"
	default:
		return "large"
	}
}

func csvPreview(headers []string, rows [][]string, hasHeader bool) string {
	start := 0
	if hasHeader {
		start = 1
	}
	limit := start + 3
	if limit > len(rows) {
		limit = len(rows)
	}
	lines := []string{"Fields: " + strings.Join(headers, ", ")}
	for _, row := range rows[start:limit] {
		lines = append(lines, strings.Join(row, " | "))
	}
	return strings.Join(lines, "\n")
}
