package config

import "errors"

type Config struct {
	ServerURL    string
	Username     string
	Password     string
	DownloadPath string
	PageSize     int // optional, defaults to 100
}

func (c *Config) Validate() error {
	if c.ServerURL == "" {
		return errors.New("server URL is required")
	}
	if c.Username == "" {
		return errors.New("username is required")
	}
	if c.Password == "" {
		return errors.New("password is required")
	}
	if c.DownloadPath == "" {
		return errors.New("download path is required")
	}
	if c.PageSize == 0 {
		c.PageSize = 100
	}
	return nil
}
