package validation

import (
	"testing"
	"time"

	"invoiceformats/pkg/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestValidator_ValidateInvoiceData_Valid(t *testing.T) {
	v := NewValidator()
	invoice := models.InvoiceData{
		Provider: models.CompanyInfo{
			ID:    uuid.New(),
			Name:  "Provider",
			Email: "provider@example.com",
			Address: models.Address{
				Street:  "123",
				City:    "City",
				Country: "US",
			},
		},
		Client: models.ClientInfo{
			ID:    uuid.New(),
			Name:  "Client",
			Email: "client@example.com",
			Address: models.Address{
				Street:  "456",
				City:    "Town",
				Country: "US",
			},
		},
		Invoice: models.InvoiceDetails{
			ID:      uuid.New(),
			Number:  "INV-001",
			Date:    time.Now(),
			DueDate: time.Now().AddDate(0, 0, 30),
			Currency: models.Currency{Code: "USD", Symbol: "$", Rate: decimal.NewFromInt(1)},
			Lines: []models.InvoiceLine{{
				ID:          uuid.New(),
				Description: "Service",
				Quantity:    decimal.NewFromInt(1),
				UnitPrice:   decimal.NewFromFloat(100),
				TaxRate:     decimal.NewFromFloat(10),
			}},
		},
	}
	invoice.Invoice.CalculateTotals()
	err := v.ValidateInvoiceData(&invoice)
	assert.NoError(t, err)
}

func TestValidator_ValidateInvoiceData_Invalid(t *testing.T) {
	v := NewValidator()
	invoice := models.InvoiceData{} // missing required fields
	err := v.ValidateInvoiceData(&invoice)
	assert.Error(t, err)
}
