// Package validation provides validation utilities for invoice data structures.
package validation

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	appErrs "invoiceformats/pkg/errors"
	"invoiceformats/pkg/i18n"
	"invoiceformats/pkg/models"

	"invoiceformats/pkg/validation/validators"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Validator wraps the go-playground validator with custom validation rules
type Validator struct {
    validate *validator.Validate
}

// NewValidator creates a new validator instance with custom rules
func NewValidator() *Validator {
    v := validator.New()
    // Register centralized validation functions
    v.RegisterValidation("currency_code", validators.CurrencyCodeValidator)
    v.RegisterValidation("iban", validators.IBANValidator)
    v.RegisterValidation("vat_id", validators.VATIDValidator)
    
    // Use JSON tag names for validation errors
    v.RegisterTagNameFunc(func(fld reflect.StructField) string {
        name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
        if name == "-" {
            return ""
        }
        return name
    })
    
    return &Validator{validate: v}
}

// ValidateInvoiceData validates the complete invoice data structure
func (v *Validator) ValidateInvoiceData(data *models.InvoiceData) error {
	// Fill in sensible defaults for missing but non-critical fields
	if data.Provider.ID == uuid.Nil {
		data.Provider.ID = uuid.New()
	}
	if data.Client.ID == uuid.Nil {
		data.Client.ID = uuid.New()
	}
	if data.Invoice.ID == uuid.Nil {
		data.Invoice.ID = uuid.New()
	}
	if data.Invoice.Number == "" {
		data.Invoice.Number = fmt.Sprintf("INV-%d", time.Now().Unix())
	}
	if data.Invoice.Language == "" {
		data.Invoice.Language = "en"
	}
	if data.Invoice.Date.IsZero() {
		data.Invoice.Date = time.Now()
	}
	if data.Invoice.DueDate.IsZero() {
		if data.Invoice.PaymentTerms.DueDays > 0 {
			data.Invoice.DueDate = data.Invoice.Date.AddDate(0, 0, data.Invoice.PaymentTerms.DueDays)
		} else {
			data.Invoice.DueDate = data.Invoice.Date.AddDate(0, 0, 30)
		}
	}
	if data.Invoice.Status == "" {
		data.Invoice.Status = models.StatusDraft
	}
	if data.Invoice.CreatedAt.IsZero() {
		data.Invoice.CreatedAt = time.Now()
	}
	if data.Invoice.UpdatedAt.IsZero() {
		data.Invoice.UpdatedAt = time.Now()
	}
	if data.Invoice.Currency.Code == "" {
		data.Invoice.Currency.Code = "EUR" // Default to EUR, or use config
		data.Invoice.Currency.Symbol = "€"
		data.Invoice.Currency.Rate = decimal.NewFromInt(1)
	}
	if data.Invoice.PaymentTerms.DueDays == 0 {
		lang := "en"
		if data.Invoice.Language != "" {
			lang = data.Invoice.Language
		}
		_ = i18n.LoadLocales(lang)
		translator := i18n.GetTranslator(lang, nil)
		data.Invoice.PaymentTerms.DueDays = 30
		data.Invoice.PaymentTerms.Description = translator("default_payment_terms")
	}
	for i := range data.Invoice.Lines {
		if data.Invoice.Lines[i].TaxRate.IsZero() {
			data.Invoice.Lines[i].TaxRate = decimal.Zero
		}
		if data.Invoice.Lines[i].Discount.IsZero() {
			data.Invoice.Lines[i].Discount = decimal.Zero
		}
	}

	if err := v.validate.Struct(data); err != nil {
		return appErrs.NewValidationError("validation failed", v.formatValidationError(err))
	}
	if err := v.validateBusinessRules(data); err != nil {
		return appErrs.NewValidationError("validation failed", err)
	}
	return nil
}

// ValidateStruct validates any struct using the validator
func (v *Validator) ValidateStruct(s interface{}) error {
	if err := v.validate.Struct(s); err != nil {
		return appErrs.NewValidationError("validation failed", v.formatValidationError(err))
	}
	return nil
}

// validateBusinessRules performs custom business logic validation
func (v *Validator) validateBusinessRules(data *models.InvoiceData) error {
    // Validate dates
    if !data.Invoice.DueDate.After(data.Invoice.Date) {
        return fmt.Errorf("due date must be after invoice date")
    }
    
    // Validate invoice lines
    if len(data.Invoice.Lines) == 0 {
        return fmt.Errorf("invoice must have at least one line item")
    }
    
    // Validate totals (this also serves as a sanity check)
    originalTotal := data.Invoice.GrandTotal
    data.Invoice.CalculateTotals()
    if !originalTotal.IsZero() && !originalTotal.Equal(data.Invoice.GrandTotal) {
        return fmt.Errorf("invoice totals are inconsistent")
    }
    
    return nil
}

// formatValidationError converts validator errors to user-friendly messages
func (v *Validator) formatValidationError(err error) error {
    var errorMessages []string
    
    for _, err := range err.(validator.ValidationErrors) {
        var message string
        
        switch err.Tag() {
        case "required":
            message = fmt.Sprintf("%s is required", err.Field())
        case "email":
            message = fmt.Sprintf("%s must be a valid email address", err.Field())
        case "min":
            message = fmt.Sprintf("%s must be at least %s", err.Field(), err.Param())
        case "max":
            message = fmt.Sprintf("%s must be at most %s", err.Field(), err.Param())
        case "len":
            message = fmt.Sprintf("%s must be exactly %s characters", err.Field(), err.Param())
        case "gt":
            message = fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param())
        case "gte":
            message = fmt.Sprintf("%s must be greater than or equal to %s", err.Field(), err.Param())
        case "lt":
            message = fmt.Sprintf("%s must be less than %s", err.Field(), err.Param())
        case "lte":
            message = fmt.Sprintf("%s must be less than or equal to %s", err.Field(), err.Param())
        case "oneof":
            message = fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param())
        case "currency_code":
            message = fmt.Sprintf("%s must be a valid ISO 4217 currency code", err.Field())
        case "iban":
            message = fmt.Sprintf("%s must be a valid IBAN", err.Field())
        case "vat_id":
            message = fmt.Sprintf("%s must be a valid VAT ID", err.Field())
        default:
            message = fmt.Sprintf("%s is invalid", err.Field())
        }
        
        errorMessages = append(errorMessages, message)
    }
    
    return fmt.Errorf("validation failed: %s", strings.Join(errorMessages, "; "))
}
