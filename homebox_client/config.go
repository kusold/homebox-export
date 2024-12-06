package homeboxclient

import (
	"net/http"
	"time"
)

type ClientConfig struct {
	BaseURL    string
	Token      string
	Timeout    time.Duration
	HTTPClient *http.Client
}

func NewClientWithConfig(config ClientConfig) (*Client, error) {
	options := []Option{
		WithToken(config.Token),
	}

	if config.HTTPClient != nil {
		options = append(options, WithHTTPClient(config.HTTPClient))
	} else if config.Timeout > 0 {
		options = append(options, WithHTTPClient(&http.Client{
			Timeout: config.Timeout,
		}))
	}

	return NewClient(config.BaseURL, options...)
}
