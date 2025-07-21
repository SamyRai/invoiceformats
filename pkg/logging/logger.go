// Package logging provides a structured, multi-level logger with environment variable control.
package logging

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// LogFields defines common structured fields for logging.
type LogFields struct {
	File         string
	Error        string
	Status       string
	URL          string
	Provider     string
	Client       string
	InvoiceNum   string
	Currency     string
	Lines        int
	EmbeddedData string
	// Extend with more fields as needed
}

// Logger is the interface for structured logging.
type Logger interface {
	Debug(msg string, fields *LogFields)
	Info(msg string, fields *LogFields)
	Warn(msg string, fields *LogFields)
	Error(msg string, fields *LogFields)
}

// zerologLogger implements Logger using zerolog.
type zerologLogger struct {
	logger  zerolog.Logger
	verbose bool
}

// NewLogger creates a new Logger instance with level and verbosity controlled by env vars.
func NewLogger() Logger {
	levelStr := strings.ToLower(os.Getenv("INVOICEGEN_LOGGING_LEVEL"))
	verbose := strings.ToLower(os.Getenv("INVOICEGEN_LOGGING_VERBOSE")) == "true"

	var level zerolog.Level
	switch levelStr {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	default:
		level = zerolog.InfoLevel
	}

	zl := zerolog.New(os.Stdout).With().Timestamp().Logger().Level(level)
	return &zerologLogger{logger: zl, verbose: verbose}
}

func (l *zerologLogger) Debug(msg string, fields *LogFields) {
	if l.verbose {
		l.logger.Debug().Fields(fieldsToMap(fields)).Msg(msg)
	} else {
		l.logger.Debug().Msg(msg)
	}
}

func (l *zerologLogger) Info(msg string, fields *LogFields) {
	if l.verbose {
		l.logger.Info().Fields(fieldsToMap(fields)).Msg(msg)
	} else {
		l.logger.Info().Msg(msg)
	}
}

func (l *zerologLogger) Warn(msg string, fields *LogFields) {
	if l.verbose {
		l.logger.Warn().Fields(fieldsToMap(fields)).Msg(msg)
	} else {
		l.logger.Warn().Msg(msg)
	}
}

func (l *zerologLogger) Error(msg string, fields *LogFields) {
	if l.verbose {
		l.logger.Error().Fields(fieldsToMap(fields)).Msg(msg)
	} else {
		l.logger.Error().Msg(msg)
	}
}

// fieldsToMap converts LogFields struct to map for zerolog.
func fieldsToMap(fields *LogFields) map[string]interface{} {
	if fields == nil {
		return nil
	}
	m := make(map[string]interface{})
	if fields.File != "" {
		m["file"] = fields.File
	}
	if fields.Error != "" {
		m["error"] = fields.Error
	}
	if fields.Status != "" {
		m["status"] = fields.Status
	}
	if fields.URL != "" {
		m["url"] = fields.URL
	}
	if fields.Provider != "" {
		m["provider"] = fields.Provider
	}
	if fields.Client != "" {
		m["client"] = fields.Client
	}
	if fields.InvoiceNum != "" {
		m["invoice_number"] = fields.InvoiceNum
	}
	if fields.Currency != "" {
		m["currency"] = fields.Currency
	}
	if fields.Lines != 0 {
		m["lines"] = fields.Lines
	}
	if fields.EmbeddedData != "" {
		m["embedded_data"] = fields.EmbeddedData
	}
	return m
}

// TODO: [HIGH] Add support for custom log outputs and structured error types. Effort: 2h. Context: Extend logger for file output and error domain integration.
