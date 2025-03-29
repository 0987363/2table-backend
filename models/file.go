package models

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	maxSize = 1024 * 1024 * 10
)

type FileType int

const (
	TypeUnknown FileType = iota
	TypePDF
	TypeImage
)

const (
	FileStatusUploading = iota
	FileStatusFinished
	FileStatusFailed
)

type File struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Status  int    `json:"status"`
	Type    string `json:"type"`
	Size    int64  `json:"size"`
	Path    string `json:"path"`
	PDFPath string `json:"pdf_path"`
}

func NewFile(db *Badger, name string) (*File, error) {
	id := GetID()
	ext := GetFileExtension(name)

	file := &File{
		ID:     id,
		Name:   name,
		Status: FileStatusUploading,
		Type:   ext,
		Path:   fmt.Sprintf("%s.%s", id, ext),
	}

	if err := db.InsertFile(FileCollection, file.ID, file); err != nil {
		return nil, err
	}

	return file, nil
}

func DetectContentType(data []byte) FileType {
	mimeType := http.DetectContentType(data)
	if strings.HasPrefix(mimeType, "image/") {
		return TypeImage
	}

	ext, _ := mime.ExtensionsByType(mimeType)
	if mimeType == "application/pdf" || (len(ext) > 0 && ext[0] == ".pdf") {
		return TypePDF
	}

	return TypeUnknown
}

func GetFileExtension(filename string) string {
	cleanName := filepath.Base(filename)
	ext := filepath.Ext(cleanName)

	if strings.HasPrefix(cleanName, ".") && len(ext) == 0 {
		return "" // .git
	}

	if len(ext) > 1 {
		return strings.ToLower(ext[1:])
	}
	return ""
}
