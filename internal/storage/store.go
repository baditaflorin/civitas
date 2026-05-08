package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/baditaflorin/civitas/internal/evidence"
)

type Store struct {
	root string
	mu   sync.RWMutex
}

func New(root string) (*Store, error) {
	if err := os.MkdirAll(filepath.Join(root, "cases"), 0o755); err != nil {
		return nil, fmt.Errorf("create storage root: %w", err)
	}
	return &Store{root: root}, nil
}

func (s *Store) CreateCase(title, description string) (evidence.Case, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	item := evidence.Case{
		ID:          newID("case"),
		Title:       title,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
		DocumentIDs: []string{},
	}
	if err := os.MkdirAll(s.caseDir(item.ID), 0o755); err != nil {
		return evidence.Case{}, fmt.Errorf("create case dir: %w", err)
	}
	if err := writeJSON(s.casePath(item.ID), item); err != nil {
		return evidence.Case{}, err
	}
	return item, nil
}

func (s *Store) ListCases() ([]evidence.Case, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries, err := os.ReadDir(filepath.Join(s.root, "cases"))
	if err != nil {
		return nil, fmt.Errorf("read cases: %w", err)
	}
	var cases []evidence.Case
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		item, err := s.readCase(entry.Name())
		if err == nil {
			cases = append(cases, item)
		}
	}
	sort.Slice(cases, func(i, j int) bool { return cases[i].CreatedAt.After(cases[j].CreatedAt) })
	return cases, nil
}

func (s *Store) GetCase(caseID string) (evidence.Case, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.readCase(caseID)
}

func (s *Store) AddDocument(caseID string, doc evidence.Document, content []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, err := s.readCase(caseID)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(s.uploadDir(caseID), 0o755); err != nil {
		return fmt.Errorf("create upload dir: %w", err)
	}
	name := filepath.Base(doc.Filename)
	if err := os.WriteFile(filepath.Join(s.uploadDir(caseID), doc.ID+"_"+name), content, 0o600); err != nil {
		return fmt.Errorf("write upload: %w", err)
	}
	if err := writeJSON(s.documentPath(caseID, doc.ID), doc); err != nil {
		return err
	}
	item.DocumentIDs = appendIfMissing(item.DocumentIDs, doc.ID)
	item.UpdatedAt = time.Now().UTC()
	return writeJSON(s.casePath(caseID), item)
}

func (s *Store) Documents(caseID string) ([]evidence.Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.documentsUnlocked(caseID)
}

func (s *Store) Search(caseID, query string) ([]evidence.SearchResult, error) {
	docs, err := s.Documents(caseID)
	if err != nil {
		return nil, err
	}
	query = stringsLowerTrim(query)
	if query == "" {
		return []evidence.SearchResult{}, nil
	}
	var results []evidence.SearchResult
	for _, doc := range docs {
		haystack := stringsLowerTrim(doc.Filename + " " + doc.Text + " " + entityText(doc.Entities))
		if !contains(haystack, query) {
			continue
		}
		results = append(results, evidence.SearchResult{
			DocumentID: doc.ID,
			Filename:   doc.Filename,
			Snippet:    snippet(doc.Text, query),
			Score:      score(haystack, query),
		})
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Score > results[j].Score })
	return results, nil
}

func (s *Store) Graph(caseID string) (evidence.Graph, error) {
	docs, err := s.Documents(caseID)
	if err != nil {
		return evidence.Graph{}, err
	}
	nodes := map[string]evidence.GraphNode{}
	edges := map[string]evidence.GraphEdge{}
	for _, doc := range docs {
		nodes[doc.ID] = evidence.GraphNode{ID: doc.ID, Type: "document", Label: doc.Filename}
		for _, entity := range doc.Entities {
			nodes[entity.ID] = evidence.GraphNode{ID: entity.ID, Type: entity.Type, Label: entity.Value}
			edgeID := doc.ID + "_" + entity.ID
			edges[edgeID] = evidence.GraphEdge{
				ID: edgeID, SourceID: doc.ID, TargetID: entity.ID, Type: "mentions", Weight: entity.Confidence,
			}
		}
	}
	return evidence.Graph{Nodes: mapNodes(nodes), Edges: mapEdges(edges)}, nil
}

func (s *Store) Timeline(caseID string) ([]evidence.TimelineEvent, error) {
	docs, err := s.Documents(caseID)
	if err != nil {
		return nil, err
	}
	var events []evidence.TimelineEvent
	for _, doc := range docs {
		events = append(events, doc.Timeline...)
	}
	sort.Slice(events, func(i, j int) bool { return events[i].When.Before(events[j].When) })
	return events, nil
}

func (s *Store) SaveExport(export evidence.Export) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, err := s.readCase(export.CaseID); err != nil {
		return err
	}
	if err := os.MkdirAll(s.exportDir(export.CaseID), 0o755); err != nil {
		return fmt.Errorf("create export dir: %w", err)
	}
	path := filepath.Join(s.exportDir(export.CaseID), export.ID+".md")
	if err := os.WriteFile(path, []byte(export.Body), 0o600); err != nil {
		return fmt.Errorf("write export: %w", err)
	}
	export.Path = path
	return writeJSON(path+".json", export)
}

func (s *Store) ReadExport(caseID, exportID string) (evidence.Export, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	path := filepath.Join(s.exportDir(caseID), filepath.Base(exportID)+".md.json")
	var item evidence.Export
	if err := readJSON(path, &item); err != nil {
		return evidence.Export{}, err
	}
	body, err := os.ReadFile(filepath.Join(s.exportDir(caseID), filepath.Base(exportID)+".md"))
	if err != nil {
		return evidence.Export{}, fmt.Errorf("read export body: %w", err)
	}
	item.Body = string(body)
	return item, nil
}

func (s *Store) documentsUnlocked(caseID string) ([]evidence.Document, error) {
	item, err := s.readCase(caseID)
	if err != nil {
		return nil, err
	}
	docs := make([]evidence.Document, 0, len(item.DocumentIDs))
	for _, docID := range item.DocumentIDs {
		var doc evidence.Document
		if err := readJSON(s.documentPath(caseID, docID), &doc); err == nil {
			docs = append(docs, doc)
		}
	}
	sort.Slice(docs, func(i, j int) bool { return docs[i].CreatedAt.Before(docs[j].CreatedAt) })
	return docs, nil
}

func (s *Store) readCase(caseID string) (evidence.Case, error) {
	var item evidence.Case
	if err := readJSON(s.casePath(caseID), &item); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return evidence.Case{}, ErrNotFound
		}
		return evidence.Case{}, err
	}
	return item, nil
}

func (s *Store) caseDir(caseID string) string {
	return filepath.Join(s.root, "cases", filepath.Base(caseID))
}

func (s *Store) casePath(caseID string) string {
	return filepath.Join(s.caseDir(caseID), "case.json")
}

func (s *Store) uploadDir(caseID string) string {
	return filepath.Join(s.caseDir(caseID), "uploads")
}

func (s *Store) exportDir(caseID string) string {
	return filepath.Join(s.caseDir(caseID), "exports")
}

func (s *Store) documentPath(caseID, docID string) string {
	return filepath.Join(s.caseDir(caseID), filepath.Base(docID)+".json")
}
