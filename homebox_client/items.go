package homeboxclient

import (
	"fmt"
	"io"
	"net/url"
	"os"
)

type ItemsService struct {
	client *Client
}

func NewItemsService(c *Client) *ItemsService {
	return &ItemsService{
		client: c,
	}
}

func (s *ItemsService) List(page, pageSize int) (*PaginationResult[Item], error) {
	u := url.Values{}
	u.Set("page", fmt.Sprintf("%d", page))
	u.Set("pageSize", fmt.Sprintf("%d", pageSize))

	req, err := s.client.newRequest("GET", "/v1/items?"+u.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var result PaginationResult[Item]
	if err := s.client.do(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *ItemsService) Get(id string) (*Item, error) {
	req, err := s.client.newRequest("GET", fmt.Sprintf("/v1/items/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var item Item
	if err := s.client.do(req, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

// func (s *ItemsService) GetAttachmentToken(itemID, attachmentID string) (*AttachmentToken, error) {
// 	req, err := s.client.newRequest("GET", fmt.Sprintf("/v1/items/%s/attachments/%s", itemID, attachmentID), nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var token AttachmentToken
// 	if err := s.client.do(req, &token); err != nil {
// 		return nil, err
// 	}

// 	return &token, nil
// }

func (s *ItemsService) DownloadAttachment(itemID, attachmentID string, destPath string) error {
	// token, err := s.GetAttachmentToken(itemID, attachmentID)
	// if err != nil {
	// 	return fmt.Errorf("failed to get attachment token: %w", err)
	// }

	// u := *s.client.baseURL
	// u.Path = filepath.Join(u.Path, "attachments", attachmentID)
	// q := u.Query()
	// q.Set("token", token.Token)
	// u.RawQuery = q.Encode()

	req, err := s.client.newRequest("GET", fmt.Sprintf("/v1/items/%s/attachments/%s", itemID, attachmentID), nil)
	if err != nil {
		return err
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download attachment: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download attachment: %s", resp.Status)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}
