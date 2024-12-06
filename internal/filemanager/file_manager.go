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

func (fm *FileManager) GenerateDirectory(item homeboxclient.Item) string {
	sanitizedName := fm.sanitizeFilename(item.Name)
	// Trim item ID to first dash
	shortId := strings.Split(item.ID, "-")[0]

	// 3M Peltor 300 Hearing Protectors_fb7115be-2ea9-4e1e-ba88-b28b3f6c0961_8b96e711-9b2c-4ebd-9fe6-a4e8ba4f1f83.pdf
	return fmt.Sprintf("%s_%s",
		sanitizedName,
		shortId)
}

func (fm *FileManager) GenerateFilename(item homeboxclient.Item, attachment homeboxclient.Attachment) string {
	ext := fm.getFileExtension(attachment.Document.Title)

	for _, a := range item.Attachments {
		if a.ID == attachment.ID {
			title := strings.TrimSuffix(attachment.Document.Title, ext)
			if title != "" {
				return fmt.Sprintf("%s%s", title, ext)
			}
			break
		}
	}
	// Fallback to attachment ID
	return fmt.Sprintf("%s%s", attachment.ID, ext)
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
