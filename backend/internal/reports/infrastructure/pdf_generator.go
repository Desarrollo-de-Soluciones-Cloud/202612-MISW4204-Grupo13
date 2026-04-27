package infrastructure

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PDFGenerator struct{}

func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

func (g *PDFGenerator) Generate(filePath string, title string, lines []string) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return err
	}

	content := buildPDFTextStream(title, lines)
	contentBytes := []byte(content)

	var out bytes.Buffer
	offsets := []int{0}

	out.WriteString("%PDF-1.4\n")
	writeObject := func(id int, body string) {
		offsets = append(offsets, out.Len())
		fmt.Fprintf(&out, "%d 0 obj\n%s\nendobj\n", id, body)
	}

	writeObject(1, "<< /Type /Catalog /Pages 2 0 R >>")
	writeObject(2, "<< /Type /Pages /Kids [3 0 R] /Count 1 >>")
	writeObject(3, "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Resources << /Font << /F1 4 0 R >> >> /Contents 5 0 R >>")
	writeObject(4, "<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>")
	writeObject(5, fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(contentBytes), content))

	xrefOffset := out.Len()
	fmt.Fprintf(&out, "xref\n0 %d\n", len(offsets))
	out.WriteString("0000000000 65535 f \n")
	for i := 1; i < len(offsets); i++ {
		fmt.Fprintf(&out, "%010d 00000 n \n", offsets[i])
	}
	fmt.Fprintf(&out, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF", len(offsets), xrefOffset)

	return os.WriteFile(filePath, out.Bytes(), 0o644)
}

func buildPDFTextStream(title string, lines []string) string {
	allLines := []string{title, ""}
	allLines = append(allLines, lines...)
	wrapped := make([]string, 0, len(allLines))
	for _, line := range allLines {
		wrapped = append(wrapped, wrapLine(line, 95)...)
	}

	var sb strings.Builder
	sb.WriteString("BT\n")
	sb.WriteString("/F1 12 Tf\n")
	sb.WriteString("50 790 Td\n")

	maxLines := 52
	if len(wrapped) > maxLines {
		wrapped = wrapped[:maxLines]
		wrapped[len(wrapped)-1] = "..."
	}

	for i, line := range wrapped {
		if i > 0 {
			sb.WriteString("T*\n")
		}
		sb.WriteString(fmt.Sprintf("(%s) Tj\n", escapePDFText(line)))
	}
	if len(wrapped) == 0 {
		sb.WriteString("( ) Tj\n")
	}

	sb.WriteString("ET")
	return sb.String()
}

func wrapLine(line string, maxLen int) []string {
	if maxLen <= 0 {
		return []string{line}
	}
	if len(line) <= maxLen {
		return []string{line}
	}

	parts := make([]string, 0)
	runes := []rune(line)
	for len(runes) > maxLen {
		parts = append(parts, string(runes[:maxLen]))
		runes = runes[maxLen:]
	}
	parts = append(parts, string(runes))
	return parts
}

func escapePDFText(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "(", "\\(")
	s = strings.ReplaceAll(s, ")", "\\)")
	return s
}
