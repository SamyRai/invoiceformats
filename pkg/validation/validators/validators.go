// Package validators provides custom field validation functions for invoice data.
package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// CurrencyCodeValidator validates ISO 4217 currency codes.
func CurrencyCodeValidator(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	if len(code) != 3 {
		return false
	}
	validCurrencies := map[string]bool{
		"USD": true, "EUR": true, "GBP": true, "JPY": true, "CHF": true,
		"CAD": true, "AUD": true, "SEK": true, "NOK": true, "DKK": true,
		"PLN": true, "CZK": true, "HUF": true, "BGN": true, "RON": true,
		"HRK": true, "RUB": true, "CNY": true, "INR": true, "BRL": true,
		"MXN": true, "ZAR": true, "KRW": true, "SGD": true, "HKD": true,
	}
	return validCurrencies[code]
}

// IBANValidator validates IBAN format.
func IBANValidator(fl validator.FieldLevel) bool {
	iban := strings.ReplaceAll(fl.Field().String(), " ", "")
	if len(iban) < 15 || len(iban) > 34 {
		return false
	}
	if len(iban) >= 4 {
		countryCode := iban[:2]
		checkDigits := iban[2:4]
		for _, c := range countryCode {
			if c < 'A' || c > 'Z' {
				return false
			}
		}
		for _, c := range checkDigits {
			if c < '0' || c > '9' {
				return false
			}
		}
		return true
	}
	return false
}

// VATIDValidator validates VAT ID format.
func VATIDValidator(fl validator.FieldLevel) bool {
	vatID := fl.Field().String()
	if len(vatID) < 4 || len(vatID) > 15 {
		return false
	}
	if len(vatID) >= 2 {
		countryCode := vatID[:2]
		for _, c := range countryCode {
			if c < 'A' || c > 'Z' {
				return false
			}
		}
		return true
	}
	return false
}

// TODO[context:validation, priority:low, effort:medium]: Add more comprehensive checks for currency, IBAN, and VAT ID formats using external libraries or official lists.
