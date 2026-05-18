package d_logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarning, "WARNING"},
		{LevelError, "ERROR"},
		{LogLevel(99), "INFO"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("LogLevel(%d).String() = %v, want %v", tt.level, got, tt.expected)
			}
		})
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
	}{
		{"DEBUG", LevelDebug},
		{"INFO", LevelInfo},
		{"WARNING", LevelWarning},
		{"ERROR", LevelError},
		{"invalid", LevelInfo},
		{"", LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ParseLogLevel(tt.input); got != tt.expected {
				t.Errorf("ParseLogLevel(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestUserLoggerManager_GetLogger(t *testing.T) {
	dir := t.TempDir()
	manager := NewUserLoggerManager(dir, LevelInfo)
	defer manager.Close()

	logger := manager.GetLogger("user1", "company1")
	if logger == nil {
		t.Fatal("GetLogger returned nil")
	}

	// Same user should return same logger
	logger2 := manager.GetLogger("user1", "company1")
	if logger != logger2 {
		t.Error("GetLogger should return same logger for same user/company")
	}

	// Different user should return different logger
	logger3 := manager.GetLogger("user2", "company1")
	if logger == logger3 {
		t.Error("GetLogger should return different logger for different user")
	}
}

func TestUserLoggerManager_WritesToFile(t *testing.T) {
	dir := t.TempDir()
	manager := NewUserLoggerManager(dir, LevelDebug)
	defer manager.Close()

	logger := manager.GetLogger("testuser", "testcompany")
	logger.Info("test message")

	// Close to flush
	manager.Close()

	logPath := filepath.Join(dir, "testuser_testcompany.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "INFO | testuser_testcompany | test message") {
		t.Errorf("log file content = %q, should contain 'INFO | testuser_testcompany | test message'", string(content))
	}
}

func TestUserLogger_LevelFiltering(t *testing.T) {
	dir := t.TempDir()
	manager := NewUserLoggerManager(dir, LevelWarning)
	defer manager.Close()

	logger := manager.GetLogger("user", "comp")
	logger.Debug("should not appear")
	logger.Info("should not appear")
	logger.Warning("should appear")
	logger.Error("should also appear")

	manager.Close()

	logPath := filepath.Join(dir, "user_comp.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	logContent := string(content)
	if strings.Contains(logContent, "DEBUG") {
		t.Error("DEBUG message should be filtered at WARNING level")
	}
	if strings.Contains(logContent, "INFO |") {
		t.Error("INFO message should be filtered at WARNING level")
	}
	if !strings.Contains(logContent, "WARNING") {
		t.Error("WARNING message should be present")
	}
	if !strings.Contains(logContent, "ERROR") {
		t.Error("ERROR message should be present")
	}
}

func TestUserLoggerManager_SetLevel(t *testing.T) {
	dir := t.TempDir()
	manager := NewUserLoggerManager(dir, LevelError)
	defer manager.Close()

	logger := manager.GetLogger("user", "comp")

	// At ERROR level, info should be filtered
	logger.Info("before level change")

	manager.SetLevel(LevelDebug)

	// After changing to DEBUG, info should appear
	logger.Info("after level change")

	manager.Close()

	logPath := filepath.Join(dir, "user_comp.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	logContent := string(content)
	if strings.Contains(logContent, "before level change") {
		t.Error("message before level change should be filtered")
	}
	if !strings.Contains(logContent, "after level change") {
		t.Error("message after level change should appear")
	}
}

func TestUserLogger_Debugf(t *testing.T) {
	dir := t.TempDir()
	manager := NewUserLoggerManager(dir, LevelDebug)
	defer manager.Close()

	logger := manager.GetLogger("u", "c")
	logger.Debugf("value=%d", 42)

	manager.Close()

	logPath := filepath.Join(dir, "u_c.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "value=42") {
		t.Errorf("Debugf formatted message not found in log: %q", string(content))
	}
}

func TestUserLogger_Close_NilFile(t *testing.T) {
	logger := &UserLogger{
		file: nil,
	}
	err := logger.Close()
	if err != nil {
		t.Errorf("Close() with nil file should return nil, got %v", err)
	}
}
