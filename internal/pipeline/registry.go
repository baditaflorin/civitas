package pipeline

import "os/exec"

type Tool struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Command     string `json:"command"`
	Available   bool   `json:"available"`
	Description string `json:"description"`
}

type Registry struct {
	tools []Tool
}

func DefaultRegistry() Registry {
	tools := []Tool{
		{Name: "Apache Tika", Kind: "document-text", Command: "tika", Description: "Office and PDF text extraction"},
		{Name: "Tesseract", Kind: "ocr", Command: "tesseract", Description: "OCR for scanned images"},
		{Name: "Pandoc", Kind: "export", Command: "pandoc", Description: "Press-ready document export"},
		{Name: "qpdf", Kind: "pdf", Command: "qpdf", Description: "PDF repair and normalization"},
		{Name: "Ghostscript", Kind: "pdf-image", Command: "gs", Description: "PDF rasterization and repair"},
		{Name: "ExifTool", Kind: "metadata", Command: "exiftool", Description: "File metadata extraction"},
		{Name: "ImageMagick", Kind: "image", Command: "magick", Description: "Image conversion and redaction prep"},
		{Name: "GraphViz", Kind: "graph", Command: "dot", Description: "Relationship graph rendering"},
		{Name: "Whisper.cpp", Kind: "transcription", Command: "whisper-cli", Description: "Local audio/video transcription"},
		{Name: "GDAL", Kind: "geospatial", Command: "gdalinfo", Description: "Geospatial evidence parsing"},
		{Name: "llama.cpp", Kind: "local-llm", Command: "llama-cli", Description: "Local summarization and Q&A"},
	}
	for i := range tools {
		_, err := exec.LookPath(tools[i].Command)
		tools[i].Available = err == nil
	}
	return Registry{tools: tools}
}

func (r Registry) Tools() []Tool {
	out := make([]Tool, len(r.tools))
	copy(out, r.tools)
	return out
}
