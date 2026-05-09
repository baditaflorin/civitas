package storage

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/baditaflorin/civitas/internal/evidence"
)

var ErrNotFound = errors.New("not found")

func newID(prefix string) string {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return prefix + "_fallback"
	}
	return prefix + "_" + hex.EncodeToString(bytes[:])
}

func writeJSON(path string, value any) error {
	body, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, body, 0o600); err != nil {
		return fmt.Errorf("write json: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("rename json: %w", err)
	}
	return nil
}

func readJSON(path string, value any) error {
	// #nosec G304 -- paths are constructed inside the storage root by Store methods.
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, value); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}
	return nil
}

func appendIfMissing(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func stringsLowerTrim(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func contains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}

func entityText(entities []evidence.Entity) string {
	values := make([]string, 0, len(entities))
	for _, entity := range entities {
		values = append(values, entity.Value)
	}
	return strings.Join(values, " ")
}

func fieldText(fields []evidence.FieldInference) string {
	values := make([]string, 0, len(fields))
	for _, field := range fields {
		values = append(values, field.Name, field.Value, field.Normalized)
	}
	return strings.Join(values, " ")
}

func snippet(text, query string) string {
	clean := strings.Join(strings.Fields(text), " ")
	lower := strings.ToLower(clean)
	idx := strings.Index(lower, query)
	if idx < 0 {
		if len(clean) > 180 {
			return clean[:180] + "..."
		}
		return clean
	}
	start := idx - 70
	if start < 0 {
		start = 0
	}
	end := idx + len(query) + 90
	if end > len(clean) {
		end = len(clean)
	}
	prefix := ""
	suffix := ""
	if start > 0 {
		prefix = "..."
	}
	if end < len(clean) {
		suffix = "..."
	}
	return prefix + clean[start:end] + suffix
}

func score(haystack, query string) float64 {
	count := strings.Count(haystack, query)
	return float64(count) + float64(len(query))/100
}

func mapNodes(nodes map[string]evidence.GraphNode) []evidence.GraphNode {
	out := make([]evidence.GraphNode, 0, len(nodes))
	for _, node := range nodes {
		out = append(out, node)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func mapEdges(edges map[string]evidence.GraphEdge) []evidence.GraphEdge {
	out := make([]evidence.GraphEdge, 0, len(edges))
	for _, edge := range edges {
		out = append(out, edge)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
