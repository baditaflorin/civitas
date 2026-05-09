package pipeline

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
)

func classifyShape(filename, contentType string, content []byte) string {
	ext := strings.ToLower(filepath.Ext(filename))
	lowerName := strings.ToLower(filename)
	switch {
	case len(content) == 0:
		return "empty"
	case ext == ".csv":
		return "csv"
	case ext == ".json" || json.Valid(content):
		if looksLikeOCDS(content) {
			return "ocds_json"
		}
		return "json"
	case ext == ".txt" || ext == ".md" || strings.HasPrefix(contentType, "text/plain"):
		return "text"
	case ext == ".html" || ext == ".htm" || strings.Contains(contentType, "html"):
		if bytes.Contains(bytes.ToLower(content), []byte("offshore leaks database")) {
			return "html_data_source"
		}
		return "html_article"
	case ext == ".pdf" || strings.Contains(contentType, "pdf"):
		return "pdf"
	case strings.HasPrefix(contentType, "image/") || ext == ".jpg" || ext == ".jpeg" || ext == ".png":
		return "image_scan"
	case strings.Contains(contentType, "ogg") || strings.Contains(contentType, "audio") || ext == ".ogg" || ext == ".mp3" || ext == ".wav":
		return "audio"
	case ext == ".zip" || strings.Contains(contentType, "zip") || lowerNameHasArchive(lowerName):
		return "archive_zip"
	default:
		return "unknown_binary"
	}
}

func looksLikeOCDS(content []byte) bool {
	var obj map[string]any
	if err := json.Unmarshal(content, &obj); err != nil {
		return false
	}
	if _, ok := obj["releases"]; ok {
		return true
	}
	_, hasOCID := obj["ocid"]
	return hasOCID
}

func lowerNameHasArchive(name string) bool {
	return strings.HasSuffix(name, ".zip") || strings.Contains(name, ".zip.")
}

func validateZip(content []byte) error {
	reader := bytes.NewReader(content)
	_, err := zip.NewReader(reader, int64(len(content)))
	return err
}
