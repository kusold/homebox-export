package downloader

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	homeboxclient "github.com/kusold/homebox-export/homebox_client"
	"github.com/kusold/homebox-export/internal/config"
)

// Mock client implementation
type mockClient struct {
	loginFunc func(username, password string) (*homeboxclient.TokenResponse, error)
}

func (m *mockClient) Login(username, password string) (*homeboxclient.TokenResponse, error) {
	if m.loginFunc != nil {
		return m.loginFunc(username, password)
	}
	return &homeboxclient.TokenResponse{Token: "test-token"}, nil
}

// Mock ItemsService implementation
type mockItemsService struct {
	listFunc               func(page, pageSize int) (*homeboxclient.PaginationResult[homeboxclient.Item], error)
	getFunc                func(id string) (*homeboxclient.Item, error)
	downloadAttachmentFunc func(itemID, attachmentID, destPath string) error
}

func (m *mockItemsService) List(page, pageSize int) (*homeboxclient.PaginationResult[homeboxclient.Item], error) {
	if m.listFunc != nil {
		return m.listFunc(page, pageSize)
	}
	return &homeboxclient.PaginationResult[homeboxclient.Item]{}, nil
}

func (m *mockItemsService) Get(id string) (*homeboxclient.Item, error) {
	if m.getFunc != nil {
		return m.getFunc(id)
	}
	return nil, nil
}

func (m *mockItemsService) DownloadAttachment(itemID, attachmentID, destPath string) error {
	if m.downloadAttachmentFunc != nil {
		return m.downloadAttachmentFunc(itemID, attachmentID, destPath)
	}
	return nil
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  config.Config
		mock    *mockClient
		wantErr bool
	}{
		{
			name: "valid configuration",
			config: config.Config{
				ServerURL:    "http://localhost",
				Username:     "test",
				Password:     "test",
				DownloadPath: t.TempDir(),
			},
			mock: &mockClient{
				loginFunc: func(username, password string) (*homeboxclient.TokenResponse, error) {
					return &homeboxclient.TokenResponse{Token: "test-token"}, nil
				},
			},
			wantErr: false,
		},
		// {
		//     name: "login failure",
		//     config: config.Config{
		//         ServerURL:    "http://localhost",
		//         Username:     "test",
		//         Password:     "test",
		//         DownloadPath: t.TempDir(),
		//     },
		//     mock: &mockClient{
		//         loginFunc: func(username, password string) (*homeboxclient.TokenResponse, error) {
		//             return nil, errors.New("login failed")
		//         },
		//     },
		//     wantErr: true,
		// },
		{
			name:   "invalid config",
			config: config.Config{
				// Missing required fields
			},
			mock:    &mockClient{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := New(tt.config, WithHomeboxClient(tt.mock))
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && d == nil {
				t.Error("New() returned nil downloader without error")
			}
		})
	}
}

func TestDownloader_DownloadAll(t *testing.T) {
	tempDir := t.TempDir()

	testItem := homeboxclient.Item{
		ID:   "test123",
		Name: "Test Item",
		Attachments: []homeboxclient.Attachment{
			{
				ID: "att123",
				Document: homeboxclient.DocumentOut{
					Title: "test.txt",
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	tests := []struct {
		name      string
		config    config.Config
		mock      *mockItemsService
		wantErr   bool
		wantFiles []string
	}{
		{
			name: "successful download",
			config: config.Config{
				ServerURL:    "http://localhost",
				Username:     "test",
				Password:     "test",
				DownloadPath: tempDir,
				PageSize:     100,
			},
			mock: &mockItemsService{
				listFunc: func(page, pageSize int) (*homeboxclient.PaginationResult[homeboxclient.Item], error) {
					if page == 1 {
						return &homeboxclient.PaginationResult[homeboxclient.Item]{
							Items: []homeboxclient.Item{testItem},
						}, nil
					}
					return &homeboxclient.PaginationResult[homeboxclient.Item]{}, nil
				},
				getFunc: func(id string) (*homeboxclient.Item, error) {
					if id == testItem.ID {
						return &testItem, nil
					}
					return nil, errors.New("item not found")
				},
				downloadAttachmentFunc: func(itemID, attachmentID, destPath string) error {
					return os.WriteFile(destPath, []byte("test content"), 0644)
				},
			},
			wantErr:   false,
			wantFiles: []string{"test.txt"},
		},
		{
			name: "list error",
			config: config.Config{
				ServerURL:    "http://localhost",
				Username:     "test",
				Password:     "test",
				DownloadPath: tempDir,
				PageSize:     100,
			},
			mock: &mockItemsService{
				listFunc: func(page, pageSize int) (*homeboxclient.PaginationResult[homeboxclient.Item], error) {
					return nil, errors.New("list error")
				},
			},
			wantErr: true,
		},
		{
			name: "download error",
			config: config.Config{
				ServerURL:    "http://localhost",
				Username:     "test",
				Password:     "test",
				DownloadPath: tempDir,
				PageSize:     100,
			},
			mock: &mockItemsService{
				listFunc: func(page, pageSize int) (*homeboxclient.PaginationResult[homeboxclient.Item], error) {
					if page == 1 {
						return &homeboxclient.PaginationResult[homeboxclient.Item]{
							Items: []homeboxclient.Item{testItem},
						}, nil
					}
					return &homeboxclient.PaginationResult[homeboxclient.Item]{}, nil
				},
				getFunc: func(id string) (*homeboxclient.Item, error) {
					return &testItem, nil
				},
				downloadAttachmentFunc: func(itemID, attachmentID, destPath string) error {
					return errors.New("download failed")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mockClient{}
			d, err := New(tt.config, WithHomeboxClient(client), WithItemService(tt.mock))
			if err != nil {
				t.Fatalf("Failed to create downloader: %v", err)
			}

			err = d.DownloadAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify expected files exist
				for _, filename := range tt.wantFiles {
					matches, err := filepath.Glob(filepath.Join(tt.config.DownloadPath, "*", filename))
					if err != nil {
						t.Errorf("Failed to glob files: %v", err)
						continue
					}
					if len(matches) == 0 {
						t.Errorf("Expected file %s not found", filename)
						continue
					}

					content, err := os.ReadFile(matches[0])
					if err != nil {
						t.Errorf("Failed to read file %s: %v", matches[0], err)
						continue
					}
					if string(content) != "test content" {
						t.Errorf("File %s content = %s, want 'test content'", matches[0], string(content))
					}
				}
			}
		})
	}
}

func TestDownloader_processItems(t *testing.T) {
	tempDir := t.TempDir()
	testItem := homeboxclient.Item{
		ID:   "test123",
		Name: "Test Item",
		Attachments: []homeboxclient.Attachment{
			{
				ID: "att123",
				Document: homeboxclient.DocumentOut{
					Title: "test.txt",
				},
			},
		},
	}

	tests := []struct {
		name    string
		items   []homeboxclient.Item
		mock    *mockItemsService
		wantErr bool
	}{
		{
			name:  "successful processing",
			items: []homeboxclient.Item{testItem},
			mock: &mockItemsService{
				getFunc: func(id string) (*homeboxclient.Item, error) {
					return &testItem, nil
				},
				downloadAttachmentFunc: func(itemID, attachmentID, destPath string) error {
					return os.WriteFile(destPath, []byte("test content"), 0644)
				},
			},
			wantErr: false,
		},
		{
			name:  "get item error",
			items: []homeboxclient.Item{testItem},
			mock: &mockItemsService{
				getFunc: func(id string) (*homeboxclient.Item, error) {
					return nil, errors.New("get error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mockClient{}
			d, err := New(config.Config{
				ServerURL:    "http://localhost",
				Username:     "user",
				Password:     "pass",
				DownloadPath: tempDir,
			}, WithHomeboxClient(client), WithItemService(tt.mock))
			if err != nil {
				t.Fatalf("Failed to create downloader: %v", err)
			}

			err = d.processItems(tt.items)
			if (err != nil) != tt.wantErr {
				t.Errorf("processItems() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
