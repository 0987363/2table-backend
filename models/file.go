package models

import (
	"context"
	"fmt"
	"image"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"gocloud.dev/blob"
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

func (file *File) ConvertPDF(storage *StorageManager) error {
	reader, err := storage.NewReader(context.Background(), file.Path, nil)
	if err != nil {
		return err
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		return err
	}

	file.PDFPath = fmt.Sprintf("%s.pdf", file.ID)
	opts := &blob.WriterOptions{ContentType: "application/pdf"}
	writer, err := storage.NewWriter(context.Background(), file.PDFPath, opts)
	if err != nil {
		return err
	}
	if err := Image2Pdf(m, writer); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}
