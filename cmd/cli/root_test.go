package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestApp_Help(t *testing.T) {
	var helpTests = []string{
		"Usage: homebox-export <command>",
		"export",
		"help",
		"version",
	}
	tests := []struct {
		name       string
		args       []string
		wantSubstr []string
	}{
		{
			name:       "no args shows help",
			args:       []string{},
			wantSubstr: helpTests,
		},
		{
			name:       "explicit help command",
			args:       []string{"help"},
			wantSubstr: helpTests,
		},
		{
			name:       "help flag",
			args:       []string{"-h"},
			wantSubstr: helpTests,
		},
		{
			name:       "long help flag",
			args:       []string{"--help"},
			wantSubstr: helpTests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			app := &App{out: &buf}
			err := app.Execute(tt.args)

			if err != nil {
				t.Errorf("Execute() error = %v", err)
				return
			}

			output := buf.String()
			for _, want := range tt.wantSubstr {
				if !strings.Contains(output, want) {
					t.Errorf("Output missing %q\nGot: %s", want, output)
				}
			}
		})
	}
}
