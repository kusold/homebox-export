package downloader

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	homeboxclient "github.com/kusold/homebox-export/homebox_client"
	"github.com/kusold/homebox-export/internal/config"
	"github.com/kusold/homebox-export/internal/filemanager"
)

type Downloader struct {
	client      HomeboxClienter
	config      config.Config
	itemService ItemServicer
	fileManager *filemanager.FileManager
}
type Option func(*Downloader)
type ItemServicer interface {
	List(page, pageSize int) (*homeboxclient.PaginationResult[homeboxclient.Item], error)
	Get(id string) (*homeboxclient.Item, error)
	DownloadAttachment(itemID, attachmentID, destPath string) error
}
type HomeboxClienter interface {
	Login(username, password string) (*homeboxclient.TokenResponse, error)
}

func New(config config.Config, options ...Option) (*Downloader, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	if err := os.MkdirAll(config.DownloadPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create download directory: %w", err)
	}

	d := &Downloader{
		config:      config,
		fileManager: filemanager.NewFileManager(config.DownloadPath),
	}
	for _, opt := range options {
		opt(d)
	}

	if d.client == nil {
		client, err := setupClient(config)
		if err != nil {
			return nil, fmt.Errorf("failed to setup client: %w", err)
		}
		d.client = client

		// I moved this inside here for mocking purposes, but agree it is odd.
		if d.itemService == nil {
			d.itemService = homeboxclient.NewItemsService(client)
		}
	}

	return d, nil
}

func WithHomeboxClient(client HomeboxClienter) Option {
	return func(d *Downloader) {
		d.client = client
	}
}
func WithItemService(is ItemServicer) Option {
	return func(d *Downloader) {
		d.itemService = is
	}
}

func (d *Downloader) DownloadAll() error {
	page := 1

	for {
		items, err := d.itemService.List(page, d.config.PageSize)
		if err != nil {
			return fmt.Errorf("failed to list items: %w", err)
		}

		if len(items.Items) == 0 {
			break
		}

		if err := d.processItems(items.Items); err != nil {
			return err
		}

		page++
	}

	return nil
}

func setupClient(config config.Config) (*homeboxclient.Client, error) {
	client, err := homeboxclient.NewClient(config.ServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	// Authenticate
	if _, err := client.Login(config.Username, config.Password); err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}
	return client, nil
}

func (d *Downloader) processItems(items []homeboxclient.Item) error {
	for _, item := range items {
		fullItem, err := d.itemService.Get(item.ID)
		if err != nil {
			return err
		}
		if err := d.processItem(*fullItem); err != nil {
			return fmt.Errorf("Error processing item %s (%s): %v", item.Name, item.ID, err)
		}
	}
	return nil
}

func (d *Downloader) processItem(item homeboxclient.Item) error {
	log.Printf("Processing item: %s (%s)", item.Name, item.ID)

	itemBytes, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item %s: %w", item.ID, err)
	} else {
		log.Printf("Item JSON: %s", string(itemBytes))
	}

	for _, attachment := range item.Attachments {
		log.Println("Processing attachment:", attachment.ID)
		subdirectory := d.fileManager.GenerateDirectory(item)
		if err := os.MkdirAll(filepath.Join(d.config.DownloadPath, subdirectory), 0755); err != nil {
			return fmt.Errorf("failed to create subdirectory: %w", err)
		}

		filename := d.fileManager.GenerateFilename(item, attachment)
		filepath := filepath.Join(d.config.DownloadPath, subdirectory, filename)

		if err := d.itemService.DownloadAttachment(item.ID, attachment.ID, filepath); err != nil {
			return fmt.Errorf("failed to download attachment %s: %w", attachment.ID, err)
		}

		log.Printf("Downloaded: %s", filename)
	}

	return nil
}
