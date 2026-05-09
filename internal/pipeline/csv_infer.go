package pipeline

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/baditaflorin/civitas/internal/evidence"
)

func sniffDelimiter(text string) rune {
	candidates := []rune{',', ';', '\t', '|'}
	firstLines := strings.Split(text, "\n")
	if len(firstLines) > 10 {
		firstLines = firstLines[:10]
	}
	best := ','
	bestScore := -1
	for _, candidate := range candidates {
		score := 0
		for _, line := range firstLines {
			score += strings.Count(line, string(candidate))
		}
		if score > bestScore {
			best = candidate
			bestScore = score
		}
	}
	return best
}

func inferHeaders(rows [][]string, fieldCount int) ([]string, bool) {
	if len(rows) == 0 || fieldCount == 0 {
		return nil, false
	}
	first := normalizeRecord(rows[0], fieldCount)
	hasHeader := true
	for _, value := range first {
		if looksNumeric(value) || strings.TrimSpace(value) == "" {
			hasHeader = false
			break
		}
	}
	headers := make([]string, fieldCount)
	for i := 0; i < fieldCount; i++ {
		if hasHeader {
			headers[i] = normalizeFieldName(first[i], i)
		} else {
			headers[i] = "column_" + strconv.Itoa(i+1)
		}
	}
	return headers, hasHeader
}

func inferCSVFields(docID string, headers []string, rows [][]string, hasHeader bool) []evidence.FieldInference {
	start := 0
	if hasHeader {
		start = 1
	}
	fields := make([]evidence.FieldInference, 0, len(headers))
	for col, header := range headers {
		values := sampleColumn(rows[start:], col)
		kind := inferKind(header, values)
		example := firstNonEmpty(values)
		fields = append(fields, newField(docID, header, kind, example, normalizeFieldValue(kind, example), 0.78, "Inferred from CSV column samples"))
	}
	return fields
}

func sampleColumn(rows [][]string, col int) []string {
	values := []string{}
	for _, row := range rows {
		if col < len(row) {
			values = append(values, strings.TrimSpace(row[col]))
		}
		if len(values) >= 50 {
			break
		}
	}
	return values
}

func normalizeRecord(row []string, fieldCount int) []string {
	out := make([]string, fieldCount)
	copy(out, row)
	return out
}

func normalizeFieldName(value string, index int) string {
	value = strings.ToLower(strings.TrimSpace(value))
	re := regexp.MustCompile(`[^a-z0-9]+`)
	value = strings.Trim(re.ReplaceAllString(value, "_"), "_")
	if value == "" {
		return "column_" + strconv.Itoa(index+1)
	}
	return value
}

func inferKind(header string, values []string) string {
	lower := strings.ToLower(header)
	if strings.Contains(lower, "email") {
		return "email"
	}
	if strings.Contains(lower, "date") {
		return "date"
	}
	if strings.Contains(lower, "amount") || strings.Contains(lower, "value") || strings.Contains(lower, "price") {
		return "money"
	}
	counts := map[string]int{}
	for _, value := range values {
		switch {
		case value == "":
		case emailPattern.MatchString(value):
			counts["email"]++
		case urlPattern.MatchString(value):
			counts["url"]++
		case moneyPattern.MatchString(value):
			counts["money"]++
		case looksNumeric(value):
			counts["number"]++
		case looksDate(value):
			counts["date"]++
		default:
			counts["text"]++
		}
	}
	best := "text"
	bestN := -1
	for kind, count := range counts {
		if count > bestN {
			best, bestN = kind, count
		}
	}
	return best
}

func normalizeFieldValue(kind, value string) string {
	if kind == "date" {
		if normalized, ok := normalizeDate(value); ok {
			return normalized
		}
	}
	return strings.TrimSpace(value)
}

func firstNonEmpty(values []string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func looksNumeric(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	_, err := strconv.ParseFloat(strings.ReplaceAll(value, ",", ""), 64)
	return err == nil
}
