package httpapi

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/baditaflorin/civitas/internal/evidence"
	"github.com/baditaflorin/civitas/internal/exporter"
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
	// #nosec G120 -- MaxBytesReader caps request bytes before multipart parsing.
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

func (a *API) getCaseState(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "case_id")
	item, err := a.deps.Store.GetCase(caseID)
	if err != nil {
		writeError(w, err, "case lookup failed")
		return
	}
	docs, err := a.deps.Store.Documents(caseID)
	if err != nil {
		writeError(w, err, "state documents failed")
		return
	}
	contents := make(map[string][]byte, len(docs))
	for _, doc := range docs {
		content, err := a.deps.Store.DocumentContent(caseID, doc.ID)
		if err != nil {
			writeError(w, err, "state source document failed")
			return
		}
		contents[doc.ID] = content
	}
	state := exporter.CaseState(item, docs, contents, a.deps.Config.Version, time.Now().UTC())
	writeJSON(w, http.StatusOK, state)
}

func (a *API) importCaseState(w http.ResponseWriter, r *http.Request) {
	var state evidence.CaseState
	if err := decodeJSON(r, &state); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Code: "bad_request", Message: "invalid case state JSON"})
		return
	}
	item, err := a.deps.Store.ImportCaseState(state)
	if err != nil {
		writeError(w, err, "import case state failed")
		return
	}
	writeJSON(w, http.StatusCreated, item)
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
	body := exporter.Markdown(item, docs, a.deps.Config.Version, a.deps.Config.CommitSHA, req.Format)
	exportID := exporter.StableExportID(caseID, req.Format, body)
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

func stableDocumentID(filename string, content []byte) string {
	hash := sha256.New()
	_, _ = hash.Write([]byte(filepath.Base(filename) + "\x00"))
	_, _ = hash.Write(content)
	return "doc_" + hex.EncodeToString(hash.Sum(nil))[:16]
}
