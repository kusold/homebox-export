package homeboxclient

// Disabled because homebox uses Swagger 2.0
// // go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml https://raw.githubusercontent.com/sysadminsmedia/homebox/refs/tags/v0.16.0/backend/app/api/static/docs/swagger.yaml

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Client struct {
	baseURL      *url.URL
	httpClient   *http.Client
	token        string
	ItemsService *ItemsService
}

type Option func(*Client)

func NewClient(baseURL string, options ...Option) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	c := &Client{
		baseURL:    parsedURL,
		httpClient: http.DefaultClient,
	}

	for _, opt := range options {
		opt(c)
	}

	return c, nil
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithToken(token string) Option {
	return func(c *Client) {
		c.token = token
	}
}

func (c *Client) newRequest(method, pathname string, body interface{}) (*http.Request, error) {
	u := *c.baseURL

	if strings.Contains(pathname, "?") {
		parts := strings.Split(pathname, "?")
		pathname = parts[0]
		if q, err := url.ParseQuery(parts[1]); err == nil {
			query := u.Query()
			for k, v := range q {
				for _, val := range v {
					query.Add(k, val)
				}
			}
			u.RawQuery = query.Encode()
		}
	}
	u.Path = path.Join(u.Path, "api", pathname)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, fmt.Errorf("failed to encode request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.token != "" {
		// req.Header.Set("Authorization", "Bearer "+c.token)
		req.Header.Set("Authorization", c.token)
	}

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
