// Package d_logger provides per-user logging capabilities for chatgraph.
// Each user gets an isolated logger that writes to a dedicated log file.
package d_logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// LogLevel represents the severity of a log message.
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarning
	LevelError
)

// String returns the string representation of a LogLevel.
func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarning:
		return "WARNING"
	case LevelError:
		return "ERROR"
	default:
		return "INFO"
	}
}

// ParseLogLevel converts a string to a LogLevel.
func ParseLogLevel(s string) LogLevel {
	switch s {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARNING":
		return LevelWarning
	case "ERROR":
		return LevelError
	default:
		return LevelInfo
	}
}

// UserLogger provides a structured logger scoped to a specific user.
type UserLogger struct {
	logger *log.Logger
	level  LogLevel
	userID string
	file   *os.File
}

// Debug logs a message at DEBUG level.
func (l *UserLogger) Debug(msg string) {
	if l.level <= LevelDebug {
		l.logger.Printf("DEBUG | %s | %s", l.userID, msg)
	}
}

// Debugf logs a formatted message at DEBUG level.
func (l *UserLogger) Debugf(format string, args ...any) {
	if l.level <= LevelDebug {
		l.logger.Printf("DEBUG | %s | %s", l.userID, fmt.Sprintf(format, args...))
	}
}

// Info logs a message at INFO level.
func (l *UserLogger) Info(msg string) {
	if l.level <= LevelInfo {
		l.logger.Printf("INFO | %s | %s", l.userID, msg)
	}
}

// Infof logs a formatted message at INFO level.
func (l *UserLogger) Infof(format string, args ...any) {
	if l.level <= LevelInfo {
		l.logger.Printf("INFO | %s | %s", l.userID, fmt.Sprintf(format, args...))
	}
}

// Warning logs a message at WARNING level.
func (l *UserLogger) Warning(msg string) {
	if l.level <= LevelWarning {
		l.logger.Printf("WARNING | %s | %s", l.userID, msg)
	}
}

// Warningf logs a formatted message at WARNING level.
func (l *UserLogger) Warningf(format string, args ...any) {
	if l.level <= LevelWarning {
		l.logger.Printf("WARNING | %s | %s", l.userID, fmt.Sprintf(format, args...))
	}
}

// Error logs a message at ERROR level.
func (l *UserLogger) Error(msg string) {
	if l.level <= LevelError {
		l.logger.Printf("ERROR | %s | %s", l.userID, msg)
	}
}

// Errorf logs a formatted message at ERROR level.
func (l *UserLogger) Errorf(format string, args ...any) {
	if l.level <= LevelError {
		l.logger.Printf("ERROR | %s | %s", l.userID, fmt.Sprintf(format, args...))
	}
}

// Close closes the underlying log file.
func (l *UserLogger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// UserLoggerManager manages per-user loggers with file-based logging.
type UserLoggerManager struct {
	mu      sync.RWMutex
	loggers map[string]*UserLogger
	dir     string
	level   LogLevel
}

// NewUserLoggerManager creates a new UserLoggerManager that stores log files
// in the specified directory. The directory is created if it does not exist.
func NewUserLoggerManager(dir string, level LogLevel) *UserLoggerManager {
	return &UserLoggerManager{
		loggers: make(map[string]*UserLogger),
		dir:     dir,
		level:   level,
	}
}

// SetLevel changes the global log level for all existing and future loggers.
func (m *UserLoggerManager) SetLevel(level LogLevel) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.level = level
	for _, l := range m.loggers {
		l.level = level
	}
}

// GetLogger returns (or creates) a logger for the given user and company.
func (m *UserLoggerManager) GetLogger(userID, companyID string) *UserLogger {
	key := userID + "_" + companyID

	m.mu.RLock()
	if l, ok := m.loggers[key]; ok {
		m.mu.RUnlock()
		return l
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring write lock
	if l, ok := m.loggers[key]; ok {
		return l
	}

	logger := m.createLogger(key)
	m.loggers[key] = logger
	return logger
}

// createLogger creates a new UserLogger. Must be called under write lock.
func (m *UserLoggerManager) createLogger(key string) *UserLogger {
	// Ensure directory exists
	if err := os.MkdirAll(m.dir, 0750); err != nil {
		// Fallback to stderr if directory creation fails
		return &UserLogger{
			logger: log.New(os.Stderr, "", log.LstdFlags),
			level:  m.level,
			userID: key,
		}
	}

	logPath := filepath.Join(m.dir, key+".log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		// Fallback to stderr
		return &UserLogger{
			logger: log.New(os.Stderr, "", log.LstdFlags),
			level:  m.level,
			userID: key,
		}
	}

	// Write to both file and stderr
	writer := io.MultiWriter(file, os.Stderr)

	return &UserLogger{
		logger: log.New(writer, "", log.LstdFlags),
		level:  m.level,
		userID: key,
		file:   file,
	}
}

// Close closes all open log files.
func (m *UserLoggerManager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, l := range m.loggers {
		l.Close()
	}
	m.loggers = make(map[string]*UserLogger)
}
