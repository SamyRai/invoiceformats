package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError_ErrorAndUnwrap(t *testing.T) {
	baseErr := errors.New("base error")
	appErr := NewValidationError("validation failed", baseErr)

	assert.Equal(t, ErrValidationFailed, appErr.Code)
	assert.Contains(t, appErr.Error(), "validation failed")
	assert.Contains(t, appErr.Error(), "base error")
	assert.Equal(t, baseErr, appErr.Unwrap())
}

func TestAppError_Constructors(t *testing.T) {
	err := NewConfigError("bad config", nil)
	assert.Equal(t, ErrConfigInvalid, err.Code)
	assert.Equal(t, "bad config", err.Message)

	nf := NewInvoiceNotFoundError("not found")
	assert.Equal(t, ErrInvoiceNotFound, nf.Code)
	assert.Equal(t, "not found", nf.Message)

	pdf := NewPDFGenerationError("pdf fail", nil)
	assert.Equal(t, ErrPDFGeneration, pdf.Code)
	assert.Equal(t, "pdf fail", pdf.Message)

	cur := NewCurrencyUnsupportedError("bad currency")
	assert.Equal(t, ErrCurrencyUnsupported, cur.Code)
	assert.Equal(t, "bad currency", cur.Message)
}
