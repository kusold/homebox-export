package logger

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{Logger: New().Logger}
	l.SetOutput(&buf)

	l.Info("test message")

	if !strings.Contains(buf.String(), "[INFO] test message") {
		t.Errorf("Info() output = %v, want %v", buf.String(), "[INFO] test message")
	}
}

func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{Logger: New().Logger}
	l.SetOutput(&buf)

	l.Error("test error")

	if !strings.Contains(buf.String(), "[ERROR] test error") {
		t.Errorf("Error() output = %v, want %v", buf.String(), "[ERROR] test error")
	}
}

func TestLogger_Debug(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{Logger: New().Logger}
	l.SetOutput(&buf)

	// Test without DEBUG env var
	l.Debug("test debug")
	if buf.String() != "" {
		t.Errorf("Debug() without DEBUG env output = %v, want empty", buf.String())
	}

	// Test with DEBUG env var
	os.Setenv("DEBUG", "1")
	defer os.Unsetenv("DEBUG")

	buf.Reset()
	l.Debug("test debug")
	if !strings.Contains(buf.String(), "[DEBUG] test debug") {
		t.Errorf("Debug() with DEBUG env output = %v, want %v", buf.String(), "[DEBUG] test debug")
	}
}

func TestLogger_Fatal(t *testing.T) {
	var buf bytes.Buffer

	var exitCode int
	osExit = func(code int) {
		exitCode = code
	}

	l := &Logger{Logger: New().Logger, exitFunc: osExit}
	l.SetOutput(&buf)

	// Override os.Exit to prevent test from exiting
	originalOsExit := osExit
	defer func() { osExit = originalOsExit }()

	l.Fatal("test fatal")

	if !strings.Contains(buf.String(), "[FATAL] test fatal") {
		t.Errorf("Fatal() output = %v, want %v", buf.String(), "[FATAL] test fatal")
	}
	if exitCode != 1 {
		t.Errorf("Fatal() exit code = %v, want 1", exitCode)
	}
}

func TestLogger_Formatted(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{Logger: New().Logger}
	l.SetOutput(&buf)

	tests := []struct {
		name    string
		logFunc func(format string, v ...interface{})
		prefix  string
		message string
		args    []interface{}
		wantLog string
	}{
		{
			name:    "Infof",
			logFunc: l.Infof,
			prefix:  "[INFO]",
			message: "test %s %d",
			args:    []interface{}{"message", 123},
			wantLog: "[INFO] test message 123",
		},
		{
			name:    "Errorf",
			logFunc: l.Errorf,
			prefix:  "[ERROR]",
			message: "test %s %d",
			args:    []interface{}{"error", 456},
			wantLog: "[ERROR] test error 456",
		},
		{
			name:    "Debugf",
			logFunc: l.Debugf,
			prefix:  "[DEBUG]",
			message: "test %s %d",
			args:    []interface{}{"debug", 789},
			wantLog: "[DEBUG] test debug 789",
		},
	}

	os.Setenv("DEBUG", "1")
	defer os.Unsetenv("DEBUG")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc(tt.message, tt.args...)
			if !strings.Contains(buf.String(), tt.wantLog) {
				t.Errorf("%s() output = %v, want %v", tt.name, buf.String(), tt.wantLog)
			}
		})
	}
}

// Mock os.Exit for testing Fatal calls
var osExit = os.Exit
