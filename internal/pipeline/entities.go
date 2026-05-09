package pipeline

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"sort"
	"strings"

	"github.com/baditaflorin/civitas/internal/evidence"
)

var (
	emailPattern   = regexp.MustCompile(`(?i)\b[A-Z0-9._%+\-]+@[A-Z0-9.\-]+\.[A-Z]{2,}\b`)
	phonePattern   = regexp.MustCompile(`\b(?:\+?\d[\d .()\-]{7,}\d)\b`)
	moneyPattern   = regexp.MustCompile(`(?i)\b(?:EUR|USD|RON|GBP|\$|€)\s?[0-9][0-9,.\s]*\b`)
	addressPattern = regexp.MustCompile(`(?im)\b\d{1,5}\s+[A-Za-z0-9 .'\-]+(?:Street|St|Road|Rd|Avenue|Ave|Boulevard|Blvd|Lane|Ln|Square|Sq|Park|Rd S)\b`)
	urlPattern     = regexp.MustCompile(`https?://[^\s"'<>]+`)
)

func extractEntities(docID, text string) []evidence.Entity {
	patterns := map[string]*regexp.Regexp{
		"address": addressPattern,
		"email":   emailPattern,
		"money":   moneyPattern,
		"phone":   phonePattern,
		"url":     urlPattern,
	}
	seen := map[string]bool{}
	out := []evidence.Entity{}
	for kind, pattern := range patterns {
		for _, match := range pattern.FindAllString(text, -1) {
			value := strings.Trim(strings.TrimSpace(match), ".,);]")
			key := kind + ":" + strings.ToLower(value)
			if value == "" || seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, evidence.Entity{ID: stableID(kind, value), Type: kind, Value: value, Confidence: 0.68, Source: docID})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Type == out[j].Type {
			return out[i].Value < out[j].Value
		}
		return out[i].Type < out[j].Type
	})
	if len(out) > 300 {
		out = out[:300]
	}
	return out
}

func stableID(prefix, value string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(value)))
	return prefix + "_" + hex.EncodeToString(sum[:])[:16]
}

func sortedFields(fields []evidence.FieldInference) []evidence.FieldInference {
	sort.SliceStable(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})
	return fields
}

func sortedAnomalies(anomalies []evidence.Anomaly) []evidence.Anomaly {
	sort.SliceStable(anomalies, func(i, j int) bool {
		return anomalies[i].Code < anomalies[j].Code
	})
	return anomalies
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}
