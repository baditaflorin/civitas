package httpapi

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/baditaflorin/civitas/internal/evidence"
	"github.com/go-chi/chi/v5"
)

type API struct {
	deps Dependencies
}

type createCaseRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type exportRequest struct {
	Format string `json:"format"`
}

func (a *API) healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) readyz(w http.ResponseWriter, _ *http.Request) {
	if _, err := a.deps.Store.ListCases(); err != nil {
		writeError(w, err, "storage is not ready")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (a *API) version(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, evidence.VersionInfo{
		Version: a.deps.Config.Version,
		Commit:  a.deps.Config.CommitSHA,
		Mode:    "github-pages-plus-docker-backend",
	})
}

func (a *API) processors(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"processors": a.deps.Pipeline.Tools()})
}

func (a *API) listCases(w http.ResponseWriter, _ *http.Request) {
	cases, err := a.deps.Store.ListCases()
	if err != nil {
		writeError(w, err, "list cases failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"cases": cases})
}

func (a *API) createCase(w http.ResponseWriter, r *http.Request) {
	var req createCaseRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Code: "bad_request", Message: "invalid JSON body"})
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Code: "bad_request", Message: "title is required"})
		return
	}
	item, err := a.deps.Store.CreateCase(req.Title, strings.TrimSpace(req.Description))
	if err != nil {
		writeError(w, err, "create case failed")
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (a *API) listDocuments(w http.ResponseWriter, r *http.Request) {
	docs, err := a.deps.Store.Documents(chi.URLParam(r, "case_id"))
	if err != nil {
		writeError(w, err, "list documents failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"documents": docs})
}

func (a *API) uploadDocument(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 100<<20)
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Code: "bad_request", Message: "multipart form required"})
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Code: "bad_request", Message: "file field is required"})
		return
	}
	defer func() {
		_ = file.Close()
	}()

	content, err := io.ReadAll(io.LimitReader(file, 100<<20))
	if err != nil {
		writeError(w, err, "read upload failed")
		return
	}
	caseID := chi.URLParam(r, "case_id")
	docID := stableDocumentID(header.Filename, content)
	doc := a.deps.Pipeline.Analyze(caseID, docID, filepath.Base(header.Filename), content)
	if err := a.deps.Store.AddDocument(caseID, doc, content); err != nil {
		writeError(w, err, "save document failed")
		return
	}
	writeJSON(w, http.StatusCreated, doc)
}

func (a *API) search(w http.ResponseWriter, r *http.Request) {
	results, err := a.deps.Store.Search(chi.URLParam(r, "case_id"), r.URL.Query().Get("q"))
	if err != nil {
		writeError(w, err, "search failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"results": results})
}

func (a *API) graph(w http.ResponseWriter, r *http.Request) {
	graph, err := a.deps.Store.Graph(chi.URLParam(r, "case_id"))
	if err != nil {
		writeError(w, err, "graph failed")
		return
	}
	writeJSON(w, http.StatusOK, graph)
}

func (a *API) timeline(w http.ResponseWriter, r *http.Request) {
	events, err := a.deps.Store.Timeline(chi.URLParam(r, "case_id"))
	if err != nil {
		writeError(w, err, "timeline failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"events": events})
}

func (a *API) debugCase(w http.ResponseWriter, r *http.Request) {
	docs, err := a.deps.Store.Documents(chi.URLParam(r, "case_id"))
	if err != nil {
		writeError(w, err, "debug case failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"documents": docs,
		"count":     len(docs),
	})
}

func (a *API) createExport(w http.ResponseWriter, r *http.Request) {
	var req exportRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Code: "bad_request", Message: "invalid JSON body"})
		return
	}
	if req.Format == "" {
		req.Format = "markdown"
	}
	caseID := chi.URLParam(r, "case_id")
	item, err := a.deps.Store.GetCase(caseID)
	if err != nil {
		writeError(w, err, "case lookup failed")
		return
	}
	docs, err := a.deps.Store.Documents(caseID)
	if err != nil {
		writeError(w, err, "export documents failed")
		return
	}
	body := buildMarkdownExport(item, docs, a.deps.Config.Version, a.deps.Config.CommitSHA, req.Format)
	exportID := stableExportID(caseID, req.Format, body)
	export := evidence.Export{
		ID:        exportID,
		CaseID:    caseID,
		Format:    req.Format,
		Body:      body,
		CreatedAt: time.Now().UTC(),
	}
	export, err = a.deps.Store.SaveExport(export)
	if err != nil {
		writeError(w, err, "save export failed")
		return
	}
	a.deps.Metrics.ExportsGenerated.Inc()
	writeJSON(w, http.StatusCreated, export)
}

func (a *API) getExport(w http.ResponseWriter, r *http.Request) {
	export, err := a.deps.Store.ReadExport(chi.URLParam(r, "case_id"), chi.URLParam(r, "export_id"))
	if err != nil {
		writeError(w, err, "read export failed")
		return
	}
	writeJSON(w, http.StatusOK, export)
}

func buildMarkdownExport(item evidence.Case, docs []evidence.Document, appVersion, commit, format string) string {
	orderedDocs := append([]evidence.Document(nil), docs...)
	sort.SliceStable(orderedDocs, func(i, j int) bool {
		return orderedDocs[i].ID < orderedDocs[j].ID
	})

	var builder strings.Builder
	builder.WriteString("# " + item.Title + "\n\n")
	builder.WriteString("Generated by Civitas safe publishing export.\n\n")
	builder.WriteString("## Export metadata\n\n")
	builder.WriteString("- schema_version: phase2.export.v1\n")
	builder.WriteString("- app_version: " + appVersion + "\n")
	builder.WriteString("- commit: " + commit + "\n")
	builder.WriteString("- format: " + format + "\n")
	builder.WriteString("- parameters: deterministic-order, redacted-contact-details, confidence-visible\n")
	builder.WriteString("- document_count: " + fmt.Sprint(len(orderedDocs)) + "\n\n")
	for _, doc := range orderedDocs {
		builder.WriteString("## " + doc.Filename + "\n\n")
		builder.WriteString("- source_id: " + doc.ID + "\n")
		builder.WriteString("- source_sha256: " + doc.SHA256 + "\n")
		builder.WriteString("- shape: " + doc.Shape + "\n")
		builder.WriteString("- state: " + doc.State + "\n")
		fmt.Fprintf(&builder, "- confidence: %.2f\n", doc.Confidence)
		builder.WriteString("- schema_version: " + doc.Provenance.SchemaVersion + "\n")
		builder.WriteString("- source_app_version: " + doc.Provenance.AppVersion + "\n")
		builder.WriteString("- processor_parameters: " + strings.Join(doc.Provenance.Parameters, ", ") + "\n\n")
		builder.WriteString(redact(doc.Summary) + "\n\n")
		if len(doc.Fields) > 0 {
			builder.WriteString("Inferred fields:\n\n")
			for _, field := range doc.Fields {
				fmt.Fprintf(&builder, "- %s (%s, %.2f): %s\n", field.Name, field.Type, field.Confidence, redact(field.Value))
			}
			builder.WriteString("\n")
		}
		if len(doc.Anomalies) > 0 {
			builder.WriteString("Warnings and next steps:\n\n")
			for _, anomaly := range doc.Anomalies {
				fmt.Fprintf(&builder, "- %s: %s Next: %s\n", anomaly.Code, anomaly.Message, anomaly.NextStep)
			}
			builder.WriteString("\n")
		}
		if len(doc.Entities) > 0 {
			builder.WriteString("Entities retained as categories, with direct contact details redacted.\n\n")
			for _, entity := range doc.Entities {
				fmt.Fprintf(&builder, "- %s: %s\n", entity.Type, redact(entity.Value))
			}
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func redact(value string) string {
	email := regexp.MustCompile(`(?i)\b[A-Z0-9._%+\-]+@[A-Z0-9.\-]+\.[A-Z]{2,}\b`)
	phone := regexp.MustCompile(`\b(?:\+?\d[\d .()\-]{7,}\d)\b`)
	value = email.ReplaceAllString(value, "[redacted-email]")
	return phone.ReplaceAllString(value, "[redacted-phone]")
}

func stableDocumentID(filename string, content []byte) string {
	hash := sha256.New()
	_, _ = hash.Write([]byte(filepath.Base(filename) + "\x00"))
	_, _ = hash.Write(content)
	return "doc_" + hex.EncodeToString(hash.Sum(nil))[:16]
}

func stableExportID(caseID, format, body string) string {
	sum := sha256.Sum256([]byte(caseID + "\x00" + format + "\x00" + body))
	return "export_" + hex.EncodeToString(sum[:])[:16]
}
