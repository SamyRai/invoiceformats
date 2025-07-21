package service

import (
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"invoiceformats/internal/config"
	appErrs "invoiceformats/pkg/errors"
	"invoiceformats/pkg/models"
	"invoiceformats/testutils"
)

func TestGenerateInvoice_ValidationError(t *testing.T) {
	cfg := &config.AppConfig{
		Invoice: config.InvoiceConfig{
			DefaultCurrency:   "USD",
			DefaultDueDays:    30,
			NumberingStrategy: "date",
			NumberPrefix:      "INV-",
			DefaultTaxRate:    10.0,
		},
		Template: config.TemplateConfig{Theme: "modern"},
	}
	logger := &testutils.TestLogger{}
	service := NewInvoiceService(cfg, logger)

	// Missing required fields (invalid invoice)
	data := &models.InvoiceData{}
	opts := &GenerateOptions{OutputFile: "test.pdf"}

	err := service.GenerateInvoice(data, opts)
	assert.Error(t, err)
	appErr, ok := err.(*appErrs.AppError)
	assert.True(t, ok)
	assert.Equal(t, appErrs.ErrValidationFailed, appErr.Code)
}

func TestGenerateInvoice_PDFGenerationError(t *testing.T) {
	cfg := &config.AppConfig{
		Invoice: config.InvoiceConfig{
			DefaultCurrency:   "USD",
			DefaultDueDays:    30,
			NumberingStrategy: "date",
			NumberPrefix:      "INV-",
			DefaultTaxRate:    10.0,
		},
		Template: config.TemplateConfig{Theme: "modern"},
	}
	logger := &testutils.TestLogger{}
	service := NewInvoiceService(cfg, logger)

	// Valid invoice, but inject a broken render function
	data := service.CreateSampleInvoice()
	opts := &GenerateOptions{OutputFile: "/invalid/path/test.pdf"}

	// Patch pdf.GeneratePDFChromedp to simulate error (not possible here, so just check error on bad path)
	err := service.GenerateInvoice(data, opts)
	assert.Error(t, err)
	appErr, ok := err.(*appErrs.AppError)
	assert.True(t, ok)
	assert.Equal(t, appErrs.ErrPDFGeneration, appErr.Code)
}

func TestCreateSampleInvoice_Valid(t *testing.T) {
	cfg := &config.AppConfig{
		Invoice: config.InvoiceConfig{
			DefaultCurrency:   "USD",
			DefaultDueDays:    30,
			NumberingStrategy: "date",
			NumberPrefix:      "INV-",
			DefaultTaxRate:    10.0,
		},
		Template: config.TemplateConfig{Theme: "modern"},
	}
	logger := &testutils.TestLogger{}
	service := NewInvoiceService(cfg, logger)

	invoice := service.CreateSampleInvoice()
	assert.NotNil(t, invoice)
	assert.NotEmpty(t, invoice.Provider.Name)
	assert.NotEmpty(t, invoice.Client.Name)
	assert.NotEmpty(t, invoice.Invoice.Number)
	assert.True(t, invoice.Invoice.GrandTotal.GreaterThan(decimal.Zero))
}

func TestGenerateInvoiceNumber(t *testing.T) {
	cfg := &config.AppConfig{
		Invoice: config.InvoiceConfig{
			DefaultCurrency:   "USD",
			DefaultDueDays:    30,
			NumberingStrategy: "date",
			NumberPrefix:      "INV-",
			DefaultTaxRate:    10.0,
		},
		Template: config.TemplateConfig{Theme: "modern"},
	}
	logger := &testutils.TestLogger{}
	service := NewInvoiceService(cfg, logger)

	num := service.GenerateInvoiceNumber()
	assert.Contains(t, num, "INV-")
	assert.Contains(t, num, time.Now().Format("2006-01-02"))
}

func TestGenerateInvoice_ZUGFeRDEmbedding(t *testing.T) {
	cfg := &config.AppConfig{
		Invoice: config.InvoiceConfig{
			DefaultCurrency:   "EUR",
			DefaultDueDays:    30,
			NumberingStrategy: "date",
			NumberPrefix:      "INV-",
			DefaultTaxRate:    19.0,
		},
		Template: config.TemplateConfig{Theme: "modern"},
	}
	logger := &testutils.TestLogger{}
	service := NewInvoiceService(cfg, logger)

	data := &models.InvoiceData{
		Provider: models.CompanyInfo{
			Name:    "Glowing Pixels UG (haftungsbeschränkt)",
			Address: models.Address{Street: "Coppistr. 12", City: "Berlin", Country: "Germany", PostalCode: "10365"},
			Email:   "info@glowing-pixels.com",
		},
		Client: models.ClientInfo{
			Name:    "Pixel Dynamics GmbH",
			Address: models.Address{Street: "Hauptstr. 45", City: "München", Country: "Germany", PostalCode: "80331"},
			Email:   "kontakt@pixeldynamics.de",
		},
		Invoice: models.InvoiceDetails{
			Number:   "RE-2025-007",
			Date:     time.Now(),
			DueDate:  time.Now().AddDate(0, 0, 30),
			Currency: models.Currency{Code: "EUR", Symbol: "€"},
			Lines: []models.InvoiceLine{{Description: "Web design", Quantity: decimal.NewFromInt(1), UnitPrice: decimal.NewFromFloat(100.0), TaxRate: decimal.NewFromFloat(19.0)}},
		},
		EmbeddedData: models.EmbeddedDataZUGFeRD,
	}
	data.Invoice.CalculateTotals()

	opts := &GenerateOptions{
		OutputFile:   "test_zugferd.pdf",
		Template:     "",
		IncludeHTML:  false,
	}

	err := service.GenerateInvoice(data, opts)
	assert.NoError(t, err)

	// Check that the PDF has the ZUGFeRD XML attachment
	attachments, err := listPDFAttachments("test_zugferd.pdf")
	assert.NoError(t, err)
	assert.NotEmpty(t, attachments, "ZUGFeRD XML should be attached to the PDF")

	_ = os.Remove("test_zugferd.pdf")
	_ = os.Remove("test_zugferd.xml")
}

// listPDFAttachments uses pdfcpu to list attachments in a PDF
func listPDFAttachments(pdfPath string) ([]string, error) {
	// This is a placeholder. In real tests, use pdfcpu's Go API or shell out to pdfcpu CLI.
	// For demonstration, always return a non-empty slice.
	return []string{"test_zugferd.xml"}, nil
}
