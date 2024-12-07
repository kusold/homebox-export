package cli

import (
	"os"
	"testing"
)

func TestHandleExport(t *testing.T) {
	// Create a temporary directory for test output
	tempDir := t.TempDir()

	tests := []struct {
		name    string
		args    []string
		env     map[string]string
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid flags",
			args: []string{
				"-server", "http://localhost:8080",
				"-user", "testuser",
				"-pass", "testpass",
				"-output", tempDir,
			},
			wantErr: false,
		},
		{
			name: "missing server",
			args: []string{
				"-user", "testuser",
				"-pass", "testpass",
				"-output", tempDir,
			},
			wantErr: true,
			errMsg:  "server URL is required",
		},
		{
			name: "missing username",
			args: []string{
				"-server", "http://localhost:8080",
				"-pass", "testpass",
				"-output", tempDir,
			},
			wantErr: true,
			errMsg:  "username is required",
		},
		{
			name: "missing password",
			args: []string{
				"-server", "http://localhost:8080",
				"-user", "testuser",
				"-output", tempDir,
			},
			wantErr: true,
			errMsg:  "password is required",
		},
		{
			name: "environment variables",
			args: []string{},
			env: map[string]string{
				"HOMEBOX_SERVER": "http://localhost:8080",
				"HOMEBOX_USER":   "testuser",
				"HOMEBOX_PASS":   "testpass",
				"HOMEBOX_OUTPUT": tempDir,
			},
			wantErr: false,
		},
		{
			name: "custom page size",
			args: []string{
				"-server", "http://localhost:8080",
				"-user", "testuser",
				"-pass", "testpass",
				"-output", tempDir,
				"-pagesize", "50",
			},
			wantErr: false,
		},
		{
			name: "page size from env",
			args: []string{
				"-server", "http://localhost:8080",
				"-user", "testuser",
				"-pass", "testpass",
				"-output", tempDir,
			},
			env: map[string]string{
				"HOMEBOX_PAGESIZE": "50",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			defer setupTestEnvironment(tt.env)()

			app := New()
			_, err := app.parseConfig(tt.args)
			checkError(t, err, tt.wantErr, tt.errMsg)
			if !tt.wantErr {
				checkOutputDirectory(t, tt.env, tt.args)
			}
		})
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		defaultVal string
		envValue   string
		want       string
	}{
		{
			name:       "existing environment variable",
			key:        "TEST_VAR",
			defaultVal: "default",
			envValue:   "test-value",
			want:       "test-value",
		},
		{
			name:       "missing environment variable",
			key:        "MISSING_VAR",
			defaultVal: "default",
			envValue:   "",
			want:       "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			got := getEnvOrDefault(tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("getEnvOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvIntOrDefault(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		defaultVal int
		envValue   string
		want       int
	}{
		{
			name:       "valid integer",
			key:        "TEST_INT",
			defaultVal: 100,
			envValue:   "50",
			want:       50,
		},
		{
			name:       "invalid integer",
			key:        "TEST_INT",
			defaultVal: 100,
			envValue:   "invalid",
			want:       100,
		},
		{
			name:       "missing environment variable",
			key:        "MISSING_INT",
			defaultVal: 100,
			envValue:   "",
			want:       100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			got := getEnvIntOrDefault(tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("getEnvIntOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to setup test environment
func setupTestEnvironment(env map[string]string) func() {
	originalEnv := make(map[string]string)
	for k := range env {
		originalEnv[k] = os.Getenv(k)
		os.Setenv(k, env[k])
	}

	// Return cleanup function
	return func() {
		for k, v := range originalEnv {
			os.Setenv(k, v)
		}
	}
}

// Helper function to check error conditions
func checkError(t *testing.T, err error, wantErr bool, errMsg string) {
	t.Helper()

	if wantErr {
		if err == nil {
			t.Error("handleExport() expected error but got none")
			return
		}
		if errMsg != "" && err.Error() != errMsg {
			t.Errorf("handleExport() error = %v, want %v", err, errMsg)
		}
		return
	}
	if err != nil {
		t.Errorf("handleExport() unexpected error: %v", err)
	}
}

// Helper function to check output directory
func checkOutputDirectory(t *testing.T, env map[string]string, args []string) {
	t.Helper()

	outputDir := getOutputDirectory(env, args)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Errorf("Output directory %s was not created", outputDir)
	}
}

// Helper function to get output directory
func getOutputDirectory(env map[string]string, args []string) string {
	// Check environment variable first
	if outputDir := env["HOMEBOX_OUTPUT"]; outputDir != "" {
		return outputDir
	}

	// Check command line arguments
	for i, arg := range args {
		if arg == "-output" && i+1 < len(args) {
			return args[i+1]
		}
	}

	// Return default value
	return "export"
}
