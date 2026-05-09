package pipeline

import (
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/baditaflorin/civitas/internal/evidence"
)

var (
	isoDatePattern   = regexp.MustCompile(`\b(20\d{2}|19\d{2})-(0[1-9]|1[0-2])-([0-2]\d|3[01])\b`)
	humanDatePattern = regexp.MustCompile(`\b(?:Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday,\s*)?(January|February|March|April|May|June|July|August|September|October|November|December)\s+([0-3]?\d),\s+(20\d{2}|19\d{2})\b`)
)

func extractTimeline(docID, text string) []evidence.TimelineEvent {
	seen := map[string]bool{}
	events := []evidence.TimelineEvent{}
	add := func(raw string) {
		normalized, ok := normalizeDate(raw)
		if !ok || seen[normalized] {
			return
		}
		seen[normalized] = true
		when, err := time.Parse("2006-01-02", normalized)
		if err != nil {
			return
		}
		events = append(events, evidence.TimelineEvent{ID: stableID("event", docID+normalized), Document: docID, When: when, Label: "Mentioned date " + normalized, Confidence: 0.76})
	}
	for _, match := range isoDatePattern.FindAllString(text, -1) {
		add(match)
	}
	for _, match := range humanDatePattern.FindAllString(text, -1) {
		add(match)
	}
	return sortedTimeline(events)
}

func timelineFromFields(docID string, fields []evidence.FieldInference) []evidence.TimelineEvent {
	events := []evidence.TimelineEvent{}
	for _, field := range fields {
		if field.Type != "date" || field.Normalized == "" {
			continue
		}
		when, err := time.Parse("2006-01-02", field.Normalized)
		if err != nil {
			continue
		}
		events = append(events, evidence.TimelineEvent{ID: stableID("event", docID+field.Name+field.Normalized), Document: docID, When: when, Label: field.Name + " " + field.Normalized, Confidence: field.Confidence})
	}
	return sortedTimeline(events)
}

func normalizeDate(value string) (string, bool) {
	value = strings.TrimSpace(value)
	if match := isoDatePattern.FindString(value); match != "" {
		return match, true
	}
	for _, layout := range []string{"Monday, January 2, 2006", "January 2, 2006", time.RFC3339} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed.Format("2006-01-02"), true
		}
	}
	if match := humanDatePattern.FindString(value); match != "" {
		clean := regexp.MustCompile(`^(Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday),\s*`).ReplaceAllString(match, "")
		if parsed, err := time.Parse("January 2, 2006", clean); err == nil {
			return parsed.Format("2006-01-02"), true
		}
	}
	return "", false
}

func looksDate(value string) bool {
	_, ok := normalizeDate(value)
	return ok
}

func sortedTimeline(events []evidence.TimelineEvent) []evidence.TimelineEvent {
	sort.SliceStable(events, func(i, j int) bool {
		if events[i].When.Equal(events[j].When) {
			return events[i].ID < events[j].ID
		}
		return events[i].When.Before(events[j].When)
	})
	return events
}
