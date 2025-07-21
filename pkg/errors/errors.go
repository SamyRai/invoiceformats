// Package errors defines domain-specific error types for InvoiceGen.
package errors

import "fmt"

// ErrorCode represents a machine-readable error code.
type ErrorCode string

const (
	ErrValidationFailed   ErrorCode = "VALIDATION_FAILED"
	ErrConfigInvalid      ErrorCode = "CONFIG_INVALID"
	ErrInvoiceNotFound    ErrorCode = "INVOICE_NOT_FOUND"
	ErrPDFGeneration      ErrorCode = "PDF_GENERATION_ERROR"
	ErrCurrencyUnsupported ErrorCode = "CURRENCY_UNSUPPORTED"
	ErrUnknown            ErrorCode = "UNKNOWN"
)

// AppError is a domain-specific error with code and message.
type AppError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

// Helper constructors
func NewValidationError(msg string, cause error) *AppError {
	return &AppError{Code: ErrValidationFailed, Message: msg, Cause: cause}
}
func NewConfigError(msg string, cause error) *AppError {
	return &AppError{Code: ErrConfigInvalid, Message: msg, Cause: cause}
}
func NewInvoiceNotFoundError(msg string) *AppError {
	return &AppError{Code: ErrInvoiceNotFound, Message: msg}
}
func NewPDFGenerationError(msg string, cause error) *AppError {
	return &AppError{Code: ErrPDFGeneration, Message: msg, Cause: cause}
}
func NewCurrencyUnsupportedError(msg string) *AppError {
	return &AppError{Code: ErrCurrencyUnsupported, Message: msg}
}
func NewAppError(code ErrorCode, msg string, cause error) *AppError {
	return &AppError{Code: code, Message: msg, Cause: cause}
}
