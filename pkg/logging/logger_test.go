package logging

import (
	"os"
	"testing"
)

func TestLoggerLevels(t *testing.T) {
	os.Setenv("INVOICEGEN_LOGGING_LEVEL", "debug")
	os.Setenv("INVOICEGEN_LOGGING_VERBOSE", "true")
	logger := NewLogger()

	logger.Debug("debug message", &LogFields{Status: "test", File: "debug.txt"})
	logger.Info("info message", &LogFields{Status: "test", File: "info.txt"})
	logger.Warn("warn message", &LogFields{Status: "test", File: "warn.txt"})
	logger.Error("error message", &LogFields{Status: "test", File: "error.txt"})
}

func TestLoggerNonVerbose(t *testing.T) {
	os.Setenv("INVOICEGEN_LOGGING_LEVEL", "info")
	os.Setenv("INVOICEGEN_LOGGING_VERBOSE", "false")
	logger := NewLogger()

	logger.Info("info message", nil)
	logger.Warn("warn message", nil)
	logger.Error("error message", nil)
}

// TODO: [MEDIUM] Add output capture and assertions for log content. Effort: 1h. Context: Improve test coverage for logger output.
