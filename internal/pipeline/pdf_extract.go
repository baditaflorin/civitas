package pipeline

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
)

// errPDFNoText is returned when the PDF parses but yields no extractable text
// (typically a scan or image-only PDF that needs OCR).
var errPDFNoText = errors.New("pdf has no extractable text layer")

// extractPDFText pulls the text layer out of a PDF using a pure-Go parser.
// It returns the concatenated plain text plus the page count. If the PDF
// parses cleanly but contains no text (e.g. scanned images), it returns
// errPDFNoText so the caller can route to the OCR processor path.
//
// The underlying parser can panic on malformed input, so this wrapper
// always recovers and turns panics into errors. Civitas accepts adversarial
// uploads — never let one take down the worker.
func extractPDFText(content []byte) (text string, pages int, err error) {
	defer func() {
		if r := recover(); r != nil {
			text = ""
			err = fmt.Errorf("pdf parser panic: %v", r)
		}
	}()
	reader, err := pdf.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return "", 0, fmt.Errorf("parse pdf: %w", err)
	}
	pages = reader.NumPage()
	plain, err := reader.GetPlainText()
	if err != nil {
		return "", pages, fmt.Errorf("extract pdf text: %w", err)
	}
	buf, err := io.ReadAll(plain)
	if err != nil {
		return "", pages, fmt.Errorf("read pdf text: %w", err)
	}
	text = strings.TrimSpace(string(buf))
	if text == "" {
		return "", pages, errPDFNoText
	}
	return text, pages, nil
}
