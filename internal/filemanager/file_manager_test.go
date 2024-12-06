package filemanager

import (
	"testing"

	homeboxclient "github.com/kusold/homebox-export/homebox_client"
)

func TestNewFileManager(t *testing.T) {
	basePath := "/test/path"
	fm := NewFileManager(basePath)
	if fm.basePath != basePath {
		t.Errorf("NewFileManager() basePath = %v, want %v", fm.basePath, basePath)
	}
}

func TestGenerateDirectory(t *testing.T) {
	tests := []struct {
		name     string
		item     homeboxclient.Item
		expected string
	}{
		{
			name: "simple item name",
			item: homeboxclient.Item{
				ID:   "abc123-456def",
				Name: "Test Item",
			},
			expected: "Test Item_abc123",
		},
		{
			name: "item with special chars",
			item: homeboxclient.Item{
				ID:   "xyz789-456def",
				Name: "Item/With:Special*Chars?",
			},
			expected: "Item_With_Special_Chars__xyz789",
		},
	}

	fm := NewFileManager("/test")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fm.GenerateDirectory(tt.item)
			if result != tt.expected {
				t.Errorf("GenerateDirectory() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGenerateFilename(t *testing.T) {
	tests := []struct {
		name       string
		item       homeboxclient.Item
		attachment homeboxclient.Attachment
		expected   string
	}{
		{
			name: "attachment with title",
			item: homeboxclient.Item{
				ID: "123",
			},
			attachment: homeboxclient.Attachment{
				ID: "att123",
				Document: homeboxclient.DocumentOut{
					Title: "document.pdf",
				},
			},
			expected: "document.pdf",
		},
		{
			name: "attachment without extension",
			item: homeboxclient.Item{ID: "123"},
			attachment: homeboxclient.Attachment{
				ID: "att123",
				Document: homeboxclient.DocumentOut{
					Title: "document",
				},
			},
			expected: "document.bin",
		},
		{
			name: "empty title",
			item: homeboxclient.Item{ID: "123"},
			attachment: homeboxclient.Attachment{
				ID: "att123",
				Document: homeboxclient.DocumentOut{
					Title: "",
				},
			},
			expected: "att123.bin",
		},
	}

	fm := NewFileManager("/test")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.item.Attachments = []homeboxclient.Attachment{tt.attachment}
			result := fm.GenerateFilename(tt.item, tt.attachment)
			if result != tt.expected {
				t.Errorf("GenerateFilename() = %v, want %v", result, tt.expected)
			}
		})
	}
}
