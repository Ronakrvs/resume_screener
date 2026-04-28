package parser

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/ledongthuc/pdf"
)

func ExtractText(filename string, data []byte) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		return extractPDF(data)
	case ".txt":
		return string(data), nil
	case ".docx":
		return extractDOCX(data)
	default:
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
}

// extractPDF tries pdftotext first (handles virtually all PDFs), falls back to
// the pure-Go reader for environments where poppler isn't installed.
func extractPDF(data []byte) (string, error) {
	if text, err := extractPDFWithPdftotext(data); err == nil && len(strings.TrimSpace(text)) > 0 {
		return cleanText(text), nil
	}
	return extractPDFGo(data)
}

func extractPDFWithPdftotext(data []byte) (string, error) {
	path, err := exec.LookPath("pdftotext")
	if err != nil {
		return "", fmt.Errorf("pdftotext not found")
	}

	// Input PDF temp file
	inFile, err := os.CreateTemp("", "resume-in-*.pdf")
	if err != nil {
		return "", err
	}
	defer os.Remove(inFile.Name())
	if _, err := inFile.Write(data); err != nil {
		inFile.Close()
		return "", err
	}
	inFile.Close()

	// Output text temp file — more reliable than stdout ("-") across versions
	outFile, err := os.CreateTemp("", "resume-out-*.txt")
	if err != nil {
		return "", err
	}
	outFile.Close()
	defer os.Remove(outFile.Name())

	// Run pdftotext; ignore non-zero exit — some PDFs produce warnings but still
	// extract text fine (exit 1 = warning, not always a hard failure)
	cmd := exec.Command(path, "-layout", "-enc", "UTF-8", inFile.Name(), outFile.Name())
	cmd.Run() // intentionally ignore error; we check the output below

	out, err := os.ReadFile(outFile.Name())
	if err != nil {
		return "", fmt.Errorf("could not read pdftotext output: %w", err)
	}
	return string(out), nil
}

func extractPDFGo(data []byte) (string, error) {
	r, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("failed to read PDF: %w", err)
	}

	var buf strings.Builder
	for i := 1; i <= r.NumPage(); i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		buf.WriteString(text)
		buf.WriteString("\n")
	}

	return cleanText(buf.String()), nil
}

func extractDOCX(data []byte) (string, error) {
	content, err := readZipEntry(bytes.NewReader(data), int64(len(data)), "word/document.xml")
	if err != nil {
		return "", fmt.Errorf("failed to read docx: %w", err)
	}
	return cleanText(stripXMLTags(content)), nil
}

func stripXMLTags(s string) string {
	var result strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			result.WriteRune(' ')
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func cleanText(s string) string {
	var result strings.Builder
	prevSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !prevSpace {
				result.WriteRune(' ')
			}
			prevSpace = true
		} else {
			result.WriteRune(r)
			prevSpace = false
		}
	}
	return strings.TrimSpace(result.String())
}
