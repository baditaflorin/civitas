package evidence

import "time"

type Case struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DocumentIDs []string  `json:"document_ids"`
}

type Document struct {
	ID          string            `json:"id"`
	CaseID      string            `json:"case_id"`
	Filename    string            `json:"filename"`
	ContentType string            `json:"content_type"`
	Size        int64             `json:"size"`
	SHA256      string            `json:"sha256"`
	Status      string            `json:"status"`
	Text        string            `json:"text,omitempty"`
	Summary     string            `json:"summary,omitempty"`
	Entities    []Entity          `json:"entities"`
	Timeline    []TimelineEvent   `json:"timeline"`
	Processors  []ProcessorStatus `json:"processors"`
	CreatedAt   time.Time         `json:"created_at"`
}

type Entity struct {
	ID         string  `json:"id"`
	Type       string  `json:"type"`
	Value      string  `json:"value"`
	Confidence float64 `json:"confidence"`
	Source     string  `json:"source"`
}

type Relationship struct {
	ID       string  `json:"id"`
	SourceID string  `json:"source_id"`
	TargetID string  `json:"target_id"`
	Type     string  `json:"type"`
	Weight   float64 `json:"weight"`
	Evidence string  `json:"evidence"`
	Document string  `json:"document_id"`
}

type TimelineEvent struct {
	ID         string    `json:"id"`
	Document   string    `json:"document_id"`
	When       time.Time `json:"when"`
	Label      string    `json:"label"`
	Confidence float64   `json:"confidence"`
}

type ProcessorStatus struct {
	Name      string    `json:"name"`
	Kind      string    `json:"kind"`
	Available bool      `json:"available"`
	Status    string    `json:"status"`
	Message   string    `json:"message,omitempty"`
	StartedAt time.Time `json:"started_at,omitempty"`
	EndedAt   time.Time `json:"ended_at,omitempty"`
}

type SearchResult struct {
	DocumentID string  `json:"document_id"`
	Filename   string  `json:"filename"`
	Snippet    string  `json:"snippet"`
	Score      float64 `json:"score"`
}

type Graph struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

type GraphNode struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Label string `json:"label"`
}

type GraphEdge struct {
	ID       string  `json:"id"`
	SourceID string  `json:"source_id"`
	TargetID string  `json:"target_id"`
	Type     string  `json:"type"`
	Weight   float64 `json:"weight"`
}

type Export struct {
	ID        string    `json:"id"`
	CaseID    string    `json:"case_id"`
	Format    string    `json:"format"`
	Path      string    `json:"path"`
	Body      string    `json:"body,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type VersionInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Mode    string `json:"mode"`
}
