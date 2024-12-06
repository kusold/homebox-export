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
	client      *homeboxclient.Client
	config      config.Config
	fileManager *filemanager.FileManager
}

func New(config config.Config) (*Downloader, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	client, err := homeboxclient.NewClient(config.ServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	if _, err := client.Login(config.Username, config.Password); err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	if err := os.MkdirAll(config.DownloadPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create download directory: %w", err)
	}

	return &Downloader{
		client:      client,
		config:      config,
		fileManager: filemanager.NewFileManager(config.DownloadPath),
	}, nil
}

func (d *Downloader) DownloadAll() error {
	page := 1

	is := homeboxclient.NewItemsService(d.client)
	for {
		items, err := is.List(page, d.config.PageSize)
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

func (d *Downloader) processItems(items []homeboxclient.Item) error {
	is := homeboxclient.NewItemsService(d.client)

	for _, item := range items {
		fullItem, err := is.Get(item.ID)
		if err != nil {
			return err
		}
		if err := d.processItem(*fullItem); err != nil {
			log.Printf("Error processing item %s (%s): %v", item.Name, item.ID, err)
			continue
		}
	}
	return nil
}

func (d *Downloader) processItem(item homeboxclient.Item) error {
	log.Printf("Processing item: %s (%s)", item.Name, item.ID)

	itemBytes, err := json.Marshal(item)
	if err != nil {
		log.Printf("Error marshaling item %s: %v", item.ID, err)
	} else {
		log.Printf("Item JSON: %s", string(itemBytes))
	}

	is := homeboxclient.NewItemsService(d.client)
	for _, attachment := range item.Attachments {
		log.Println("Processing attachment:", attachment.ID)
		subdirectory := d.fileManager.GenerateDirectory(item)
		if err := os.MkdirAll(filepath.Join(d.config.DownloadPath, subdirectory), 0755); err != nil {
			return fmt.Errorf("failed to create subdirectory: %w", err)
		}

		filename := d.fileManager.GenerateFilename(item, attachment)
		filepath := filepath.Join(d.config.DownloadPath, subdirectory, filename)

		if err := is.DownloadAttachment(item.ID, attachment.ID, filepath); err != nil {
			return fmt.Errorf("failed to download attachment %s: %w", attachment.ID, err)
		}

		log.Printf("Downloaded: %s", filename)
	}

	return nil
}
