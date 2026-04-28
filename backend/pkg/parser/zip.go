package parser

import (
	"archive/zip"
	"fmt"
	"io"
)

type zipEntry struct {
	name    string
	content string
}

func newZipReader(r io.ReaderAt, size int64) ([]zipEntry, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	var entries []zipEntry
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			continue
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			continue
		}
		entries = append(entries, zipEntry{name: f.Name, content: string(data)})
	}
	return entries, nil
}

func readZipEntry(r io.ReaderAt, size int64, filename string) (string, error) {
	entries, err := newZipReader(r, size)
	if err != nil {
		return "", err
	}
	for _, f := range entries {
		if f.name == filename {
			return f.content, nil
		}
	}
	return "", fmt.Errorf("entry %s not found", filename)
}
