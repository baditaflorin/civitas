package pipeline

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func normalizeBytes(content []byte) string {
	if len(content) >= 3 && content[0] == 0xef && content[1] == 0xbb && content[2] == 0xbf {
		content = content[3:]
	}
	text := string(content)
	if !utf8.ValidString(text) {
		text = strings.ToValidUTF8(text, "�")
	}
	return normalizeText(text)
}

func normalizeText(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	text = strings.ReplaceAll(text, "\u00a0", " ")
	return strings.TrimSpace(text)
}

func preview(text string) string {
	return previewLong(text, 320)
}

func previewLong(text string, max int) string {
	clean := strings.Join(strings.FieldsFunc(text, func(r rune) bool {
		return unicode.IsSpace(r)
	}), " ")
	if len(clean) > max {
		return clean[:max] + "..."
	}
	return clean
}
