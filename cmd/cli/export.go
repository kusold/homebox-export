package cli

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/kusold/homebox-export/internal/config"
	"github.com/kusold/homebox-export/internal/downloader"
)

func (a *App) handleExport(args []string) error {
	cmd := flag.NewFlagSet("export", flag.ExitOnError)

	var config config.Config

	// Default to environment variables if available
	cmd.StringVar(&config.ServerURL, "server", os.Getenv("HOMEBOX_SERVER"), "Homebox server URL (required)")
	cmd.StringVar(&config.Username, "user", os.Getenv("HOMEBOX_USER"), "Username for authentication (required)")
	cmd.StringVar(&config.Password, "pass", os.Getenv("HOMEBOX_PASS"), "Password for authentication (required)")
	cmd.StringVar(&config.DownloadPath, "output", getEnvOrDefault("HOMEBOX_OUTPUT", "downloads"), "Output directory")
	cmd.IntVar(&config.PageSize, "pagesize", getEnvIntOrDefault("HOMEBOX_PAGESIZE", 100), "Number of items per page")

	if err := cmd.Parse(args); err != nil {
		return err
	}

	// Validate required flags
	if config.ServerURL == "" {
		return fmt.Errorf("server URL is required")
	}
	if config.Username == "" {
		return fmt.Errorf("username is required")
	}
	if config.Password == "" {
		return fmt.Errorf("password is required")
	}

	d, err := downloader.New(config)
	if err != nil {
		return fmt.Errorf("failed to initialize downloader: %w", err)
	}

	return d.DownloadAll()
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}
