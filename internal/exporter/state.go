package exporter

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"sort"
	"time"

	"github.com/baditaflorin/civitas/internal/evidence"
)

const CaseStateSchemaVersion = "civitas.case_state.v1"

func CaseState(item evidence.Case, docs []evidence.Document, contents map[string][]byte, appVersion string, exportedAt time.Time) evidence.CaseState {
	orderedDocs := append([]evidence.Document(nil), docs...)
	sort.SliceStable(orderedDocs, func(i, j int) bool {
		return orderedDocs[i].ID < orderedDocs[j].ID
	})
	stateDocs := make([]evidence.StateDocument, 0, len(orderedDocs))
	for _, doc := range orderedDocs {
		content := contents[doc.ID]
		sum := sha256.Sum256(content)
		stateDocs = append(stateDocs, evidence.StateDocument{
			Document:      doc,
			ContentBase64: base64.StdEncoding.EncodeToString(content),
			ContentSHA256: hex.EncodeToString(sum[:]),
		})
	}
	return evidence.CaseState{
		SchemaVersion: CaseStateSchemaVersion,
		AppVersion:    appVersion,
		ExportedAt:    exportedAt,
		Case:          item,
		Documents:     stateDocs,
	}
}
