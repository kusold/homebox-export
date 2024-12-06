package filemanager

import (
	"fmt"
	"path/filepath"
	"strings"

	homeboxclient "github.com/kusold/homebox-export/homebox_client"
)

type FileManager struct {
	basePath string
}

func NewFileManager(basePath string) *FileManager {
	return &FileManager{
		basePath: basePath,
	}
}

func (fm *FileManager) GenerateFilename(item homeboxclient.Item, attachment homeboxclient.Attachment) string {
	sanitizedName := fm.sanitizeFilename(item.Name)
	ext := fm.getFileExtension(attachment.Document.Title)

	return fmt.Sprintf("%s_%s_%s%s",
		sanitizedName,
		item.ID,
		attachment.ID,
		ext)
}

func (fm *FileManager) sanitizeFilename(filename string) string {
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := filename

	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}

	result = strings.TrimSpace(result)
	if len(result) > 50 {
		result = result[:50]
	}

	return result
}

func (fm *FileManager) getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ".bin"
	}
	return ext
}
