package zugferd_test

import (
	"testing"
	"time"

	"invoiceformats/pkg/models"
	"invoiceformats/providers/zugferd"

	"github.com/shopspring/decimal"
)

func TestBuildBasicXML_ValidInvoice(t *testing.T) {
	inv := models.InvoiceData{
		Provider: models.CompanyInfo{
			Name: "Test Seller",
			VATID: "DE123456789",
			Address: models.Address{
				Street:     "Teststr. 1",
				City:       "Berlin",
				PostalCode: "10115",
				Country:    "DE",
			},
		},
		Client: models.ClientInfo{
			Name: "Test Buyer",
			VATID: "DE987654321",
			Address: models.Address{
				Street:     "Kaufstr. 2",
				City:       "Munich",
				PostalCode: "80331",
				Country:    "DE",
			},
		},
		Invoice: models.InvoiceDetails{
			Number: "INV-001",
			Date:   time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
			GrandTotal: decimal.NewFromFloat(100.00),
			Currency: models.Currency{Code: "EUR"},
			Lines: []models.InvoiceLine{
				{
					Description: "Service",
					Quantity:    decimal.NewFromFloat(1),
					UnitPrice:   decimal.NewFromFloat(100.00),
					Total:       decimal.NewFromFloat(100.00),
					TaxRate:     decimal.NewFromFloat(19.0),
					TaxAmount:   decimal.NewFromFloat(19.00),
				},
			},
		},
	}
	builder := zugferd.ZUGFeRDBasicXMLBuilder{}
	xmlData, err := builder.BuildXML(inv)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(xmlData) == 0 {
		t.Error("expected non-empty XML output")
	}
}

func TestBuildBasicXML_MissingFields(t *testing.T) {
	inv := models.InvoiceData{}
	builder := zugferd.ZUGFeRDBasicXMLBuilder{}
	_, err := builder.BuildXML(inv)
	if err == nil {
		t.Error("expected error for missing fields, got nil")
	}
}

func TestBuildBasicXML_InvalidDate(t *testing.T) {
	inv := models.InvoiceData{
		Provider: models.CompanyInfo{Name: "Seller"},
		Client: models.ClientInfo{Name: "Buyer"},
		Invoice: models.InvoiceDetails{
			Number: "INV-001",
			Date:   time.Time{}, // Invalid date
			GrandTotal: decimal.NewFromFloat(100.00),
			Currency: models.Currency{Code: "EUR"},
		},
	}
	builder := zugferd.ZUGFeRDBasicXMLBuilder{}
	_, err := builder.BuildXML(inv)
	if err == nil {
		t.Error("expected error for invalid IssueDate, got nil")
	}
}

// TODO: Add more tests for line items, taxes, and edge cases as model expands
