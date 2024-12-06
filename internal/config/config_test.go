package config

import (
	"testing"
)

func TestConfig_Validate(t *testing.T) {
    tests := []struct {
        name          string
        config        Config
        wantErr      bool
        wantPageSize int // Add field to check expected PageSize
    }{
        {
            name: "valid config",
            config: Config{
                ServerURL:    "http://localhost:8080",
                Username:     "user",
                Password:     "pass",
                DownloadPath: "/tmp",
                PageSize:     100,
            },
            wantErr:      false,
            wantPageSize: 100,
        },
        {
            name: "zero page size gets default",
            config: Config{
                ServerURL:    "http://localhost:8080",
                Username:     "user",
                Password:     "pass",
                DownloadPath: "/tmp",
                PageSize:     0,
            },
            wantErr:      false,
            wantPageSize: 100, // Should be set to default
        },
        {
            name: "missing server URL",
            config: Config{
                Username:     "user",
                Password:     "pass",
                DownloadPath: "/tmp",
            },
            wantErr:      true,
            wantPageSize: 100,
        },
        {
            name: "missing username",
            config: Config{
                ServerURL:    "http://localhost:8080",
                Password:     "pass",
                DownloadPath: "/tmp",
            },
            wantErr:      true,
            wantPageSize: 100,
        },
        {
            name: "missing password",
            config: Config{
                ServerURL:    "http://localhost:8080",
                Username:     "user",
                DownloadPath: "/tmp",
            },
            wantErr:      true,
            wantPageSize: 100,
        },
        {
            name: "missing download path",
            config: Config{
                ServerURL: "http://localhost:8080",
                Username: "user",
                Password: "pass",
            },
            wantErr:      true,
            wantPageSize: 100,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.config.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }

            // Check if PageSize was set to default value when appropriate
            if err == nil && tt.config.PageSize != tt.wantPageSize {
                t.Errorf("Validate() PageSize = %v, want %v", tt.config.PageSize, tt.wantPageSize)
            }
        })
    }
}
