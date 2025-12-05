// Package core provides structured logging functionality for the IDE.
package core

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// Performance: Pool of strings.Builder for logger
var (
	loggerBuilderPool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
)

// LogLevel represents logging levels
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// ParseLogLevel parses a string into a LogLevel
func ParseLogLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARN", "WARNING":
		return LevelWarn
	case "ERROR":
		return LevelError
	default:
		return LevelInfo
	}
}

// StandardLogger implements the Logger interface using Go's standard log package
type StandardLogger struct {
	logger   *log.Logger
	minLevel LogLevel
}

// NewStandardLogger creates a new standard logger
func NewStandardLogger(w io.Writer, minLevel LogLevel) *StandardLogger {
	if w == nil {
		w = os.Stderr
	}

	return &StandardLogger{
		logger:   log.New(w, "", 0), // No prefix, we'll format ourselves
		minLevel: minLevel,
	}
}

// Debug logs a debug message
func (l *StandardLogger) Debug(msg string, fields ...Field) {
	if l.minLevel <= LevelDebug {
		l.log(LevelDebug, msg, fields...)
	}
}

// Info logs an info message
func (l *StandardLogger) Info(msg string, fields ...Field) {
	if l.minLevel <= LevelInfo {
		l.log(LevelInfo, msg, fields...)
	}
}

// Warn logs a warning message
func (l *StandardLogger) Warn(msg string, fields ...Field) {
	if l.minLevel <= LevelWarn {
		l.log(LevelWarn, msg, fields...)
	}
}

// Error logs an error message
func (l *StandardLogger) Error(msg string, fields ...Field) {
	if l.minLevel <= LevelError {
		l.log(LevelError, msg, fields...)
	}
}

func (l *StandardLogger) log(level LogLevel, msg string, fields ...Field) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Use pooled builder to avoid allocations
	fieldStr := loggerBuilderPool.Get().(*strings.Builder)
	fieldStr.Reset()
	defer loggerBuilderPool.Put(fieldStr)
	
	if len(fields) > 0 {
		fieldStr.WriteString(" [")
		for i, field := range fields {
			if i > 0 {
				fieldStr.WriteString(", ")
			}
			fieldStr.WriteString(field.Key)
			fieldStr.WriteString("=")
			fieldStr.WriteString(fmt.Sprintf("%v", field.Value))
		}
		fieldStr.WriteString("]")
	}

	// Use pooled builder for final message too
	msgBuilder := loggerBuilderPool.Get().(*strings.Builder)
	msgBuilder.Reset()
	defer loggerBuilderPool.Put(msgBuilder)
	
	msgBuilder.WriteString(timestamp)
	msgBuilder.WriteString(" [")
	msgBuilder.WriteString(level.String())
	msgBuilder.WriteString("] ")
	msgBuilder.WriteString(msg)
	msgBuilder.WriteString(fieldStr.String())

	l.logger.Println(msgBuilder.String())
}

// NoopLogger is a logger that does nothing
type NoopLogger struct{}

// NewNoopLogger creates a new no-op logger
func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

// Debug does nothing
func (l *NoopLogger) Debug(msg string, fields ...Field) {}

// Info does nothing
func (l *NoopLogger) Info(msg string, fields ...Field) {}

// Warn does nothing
func (l *NoopLogger) Warn(msg string, fields ...Field) {}

// Error does nothing
func (l *NoopLogger) Error(msg string, fields ...Field) {}
