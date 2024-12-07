package cli

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func TestVersionInfo(t *testing.T) {
	// Save original values to restore after test
	origVersion := version
	origCommit := commit
	origDate := date
	defer func() {
		version = origVersion
		commit = origCommit
		date = origDate
	}()

	tests := []struct {
		name    string
		version string
		commit  string
		date    string
		want    []string
	}{
		{
			name:    "default values",
			version: "dev",
			commit:  "none",
			date:    "unknown",
			want: []string{
				"homebox-export dev",
				"(none)",
				"built unknown",
				runtime.GOOS + "/" + runtime.GOARCH,
			},
		},
		{
			name:    "specific version",
			version: "v1.0.0",
			commit:  "abc123",
			date:    "2024-01-01",
			want: []string{
				"homebox-export v1.0.0",
				"(abc123)",
				"built 2024-01-01",
				runtime.GOOS + "/" + runtime.GOARCH,
			},
		},
		{
			name:    "long commit hash",
			version: "v1.0.0",
			commit:  "abc123def456789",
			date:    "2024-01-01",
			want: []string{
				"homebox-export v1.0.0",
				"(abc123def456789)",
				"built 2024-01-01",
				runtime.GOOS + "/" + runtime.GOARCH,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version = tt.version
			commit = tt.commit
			date = tt.date

			got := versionInfo()

			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf("versionInfo() = %q, want it to contain %q", got, want)
				}
			}

			expectedFormat := "homebox-export %s (%s) - built %s\n%s/%s"
			expected := fmt.Sprintf(expectedFormat,
				tt.version,
				tt.commit,
				tt.date,
				runtime.GOOS,
				runtime.GOARCH,
			)
			if got != expected {
				t.Errorf("versionInfo() = %q, want %q", got, expected)
			}
		})
	}
}

func TestVersionInfo_Components(t *testing.T) {
	info := versionInfo()
	components := []string{
		"homebox-export",
		version,
		"(",
		commit,
		")",
		"built",
		date,
		runtime.GOOS,
		runtime.GOARCH,
	}

	lastIndex := -1
	for _, component := range components {
		index := strings.Index(info, component)
		if index == -1 {
			t.Errorf("versionInfo() missing component %q", component)
			continue
		}
		if index <= lastIndex {
			t.Errorf("versionInfo() component %q in wrong position", component)
		}
		lastIndex = index
	}
}

func TestVersionInfo_Format(t *testing.T) {
	expectedFormat := "homebox-export %s (%s) - built %s\n%s/%s"
	info := versionInfo()
	expected := fmt.Sprintf(expectedFormat,
		version,
		commit,
		date,
		runtime.GOOS,
		runtime.GOARCH,
	)

	if info != expected {
		t.Errorf("versionInfo() format mismatch\ngot:  %q\nwant: %q", info, expected)
	}
}
