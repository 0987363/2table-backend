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

type File struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Path string `json:"path"`
	Size int    `json:"size"`
}

func NewFile(name string, data []byte) *File {
	id := GetID()
	ext := GetFileExtension(name)

	file := File{
		ID:   id,
		Name: name,
		Type: ext,
		Path: fmt.Sprintf("%s.%s", id, ext),
		Size: len(data),
	}

	return &file
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
