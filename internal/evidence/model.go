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
	Shape       string            `json:"shape"`
	State       string            `json:"state"`
	Confidence  float64           `json:"confidence"`
	Text        string            `json:"text,omitempty"`
	Preview     string            `json:"preview,omitempty"`
	Summary     string            `json:"summary,omitempty"`
	Fields      []FieldInference  `json:"fields"`
	Anomalies   []Anomaly         `json:"anomalies"`
	Entities    []Entity          `json:"entities"`
	Timeline    []TimelineEvent   `json:"timeline"`
	Processors  []ProcessorStatus `json:"processors"`
	Provenance  Provenance        `json:"provenance"`
	Parse       ParseMetrics      `json:"parse"`
	CreatedAt   time.Time         `json:"created_at"`
}

type FieldInference struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Value      string  `json:"value"`
	Normalized string  `json:"normalized,omitempty"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason"`
}

type Anomaly struct {
	ID         string  `json:"id"`
	Code       string  `json:"code"`
	Severity   string  `json:"severity"`
	Message    string  `json:"message"`
	Why        string  `json:"why"`
	NextStep   string  `json:"next_step"`
	Confidence float64 `json:"confidence"`
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

type Provenance struct {
	SchemaVersion string   `json:"schema_version"`
	AppVersion    string   `json:"app_version"`
	SourceID      string   `json:"source_id"`
	SourceName    string   `json:"source_name"`
	SourceSHA256  string   `json:"source_sha256"`
	Parameters    []string `json:"parameters"`
}

type ParseMetrics struct {
	DurationMS int64  `json:"duration_ms"`
	SizeBucket string `json:"size_bucket"`
	RowCount   int    `json:"row_count,omitempty"`
	FieldCount int    `json:"field_count,omitempty"`
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

type CaseState struct {
	SchemaVersion string          `json:"schema_version"`
	AppVersion    string          `json:"app_version"`
	ExportedAt    time.Time       `json:"exported_at"`
	Case          Case            `json:"case"`
	Documents     []StateDocument `json:"documents"`
}

type StateDocument struct {
	Document      Document `json:"document"`
	ContentBase64 string   `json:"content_base64"`
	ContentSHA256 string   `json:"content_sha256"`
}

type VersionInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Mode    string `json:"mode"`
}
