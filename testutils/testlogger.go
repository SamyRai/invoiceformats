package testutils

import "invoiceformats/pkg/logging"

// TestLogger is a stub logger for use in tests.
type TestLogger struct{}

func (t *TestLogger) Debug(msg string, fields *logging.LogFields)   {}
func (t *TestLogger) Info(msg string, fields *logging.LogFields)    {}
func (t *TestLogger) Warn(msg string, fields *logging.LogFields)    {}
func (t *TestLogger) Error(msg string, fields *logging.LogFields)   {}
