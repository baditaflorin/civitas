package pipeline

import (
	"fmt"
	"strconv"

	"github.com/baditaflorin/civitas/internal/evidence"
)

func inferJSONFields(docID string, obj any, shape string) []evidence.FieldInference {
	if shape != "ocds_json" {
		return genericJSONFields(docID, obj)
	}
	root, ok := obj.(map[string]any)
	if !ok {
		return genericJSONFields(docID, obj)
	}
	releases, _ := root["releases"].([]any)
	if len(releases) == 0 {
		return genericJSONFields(docID, obj)
	}
	release, _ := releases[0].(map[string]any)
	fields := []evidence.FieldInference{}
	addPath := func(name, kind string, path ...string) {
		if value, ok := stringPath(release, path...); ok {
			fields = append(fields, newField(docID, name, kind, value, normalizeFieldValue(kind, value), 0.92, "Detected from OCDS path "+fmt.Sprint(path)))
		}
	}
	addPath("ocid", "id", "ocid")
	addPath("buyer", "organization", "buyer", "name")
	addPath("tender_title", "text", "tender", "title")
	addPath("tender_value", "money", "tender", "value", "amount")
	addPath("currency", "currency", "tender", "value", "currency")
	addPath("release_date", "date", "date")
	return fields
}

func genericJSONFields(docID string, obj any) []evidence.FieldInference {
	switch typed := obj.(type) {
	case map[string]any:
		fields := make([]evidence.FieldInference, 0, len(typed))
		for key, value := range typed {
			text := fmt.Sprint(value)
			if len(text) > 120 {
				text = text[:120] + "..."
			}
			fields = append(fields, newField(docID, normalizeFieldName(key, 0), inferKind(key, []string{text}), text, text, 0.72, "Detected from top-level JSON key"))
			if len(fields) >= 20 {
				break
			}
		}
		return fields
	default:
		return nil
	}
}

func stringPath(root map[string]any, path ...string) (string, bool) {
	var current any = root
	for _, part := range path {
		m, ok := current.(map[string]any)
		if !ok {
			return "", false
		}
		current, ok = m[part]
		if !ok {
			return "", false
		}
	}
	switch value := current.(type) {
	case string:
		return value, true
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64), true
	default:
		return fmt.Sprint(value), true
	}
}
