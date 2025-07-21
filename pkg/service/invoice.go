// Package service provides business logic for invoice operations.
package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"invoiceformats/internal/config"
	"invoiceformats/pkg/compliance"
	"invoiceformats/pkg/di"
	appErrs "invoiceformats/pkg/errors"
	"invoiceformats/pkg/i18n"
	"invoiceformats/pkg/interfaces"
	"invoiceformats/pkg/logging"
	"invoiceformats/pkg/models"
	"invoiceformats/pkg/pdf"
	"invoiceformats/pkg/render"
	"invoiceformats/pkg/validation"
	"invoiceformats/pkg/xml"
)

// InvoiceService handles invoice generation and related operations
// Inject ZUGFeRDInvoiceXMLBuilder for DI

type InvoiceService struct {
	config    *config.AppConfig
	logger    logging.Logger
	validator *validation.Validator
}

// NewInvoiceService creates a new invoice service instance
func NewInvoiceService(cfg *config.AppConfig, logger logging.Logger) *InvoiceService {
	return &InvoiceService{
		config:    cfg,
		logger:    logger,
		validator: validation.NewValidator(),
	}
}

// GenerateOptions holds options for invoice generation
type GenerateOptions struct {
	OutputFile     string
	Template       string
	Currency       string
	TaxRate        *float64
	IncludeHTML    bool
	DryRun         bool
	ValidateOnly   bool
	EnableZUGFeRD  bool
	EmbeddedDataProvider interfaces.PDFEmbeddedDataProvider
}

// GenerateInvoice creates an invoice PDF from the provided data
func (s *InvoiceService) GenerateInvoice(data *models.InvoiceData, opts *GenerateOptions) error {
	s.logger.Info("Starting invoice generation", &logging.LogFields{
		InvoiceNum: data.Invoice.Number,
		File: opts.OutputFile,
	})

	// Debug: Log currency before applying defaults
	s.logger.Debug("Currency before applying defaults", &logging.LogFields{
		Currency: data.Invoice.Currency.Code,
		Status: data.Invoice.Currency.Symbol,
	})

	// Apply service defaults FIRST
	s.applyDefaults(data, opts)

	// Debug: Log currency after applying defaults
	s.logger.Debug("Currency after applying defaults", &logging.LogFields{
		Currency: data.Invoice.Currency.Code,
		Status: data.Invoice.Currency.Symbol,
	})

	// Validate the invoice data AFTER applying defaults
	if err := s.validator.ValidateInvoiceData(data); err != nil {
		s.logger.Error("Invoice validation failed", &logging.LogFields{Error: err.Error()})
		return appErrs.NewValidationError("validation failed", err)
	}

	if opts.ValidateOnly {
		s.logger.Info("Validation successful, skipping generation (validate-only mode)", nil)
		return nil
	}

	// Calculate totals
	data.Invoice.CalculateTotals()

	// Auto-set recurrence if all lines have the same period or only one line with period
	periodSet := make(map[string]struct{})
	for _, line := range data.Invoice.Lines {
		if line.Period != "" {
			periodSet[line.Period] = struct{}{}
		}
	}
	if len(periodSet) == 1 && len(data.Invoice.Lines) > 0 {
		for p := range periodSet {
			data.Invoice.Recurrence = p
		}
	}

	// Update service to use new RenderHTML signature and set EmbeddedData
	// Render HTML and extract embedded data type from template
	html, err := render.RenderHTML(*data, opts.Template)
	if err != nil {
		s.logger.Error("HTML rendering failed", &logging.LogFields{Error: err.Error()})
		return appErrs.NewPDFGenerationError("failed to render HTML", err)
	}

	if opts.DryRun {
		s.logger.Info("Dry run mode - would generate PDF", &logging.LogFields{File: opts.OutputFile})
		return nil
	}

	// Save HTML if requested
	if opts.IncludeHTML {
		htmlFile := changeExtension(opts.OutputFile, ".html")
		if err := os.WriteFile(htmlFile, []byte(html), 0644); err != nil {
			s.logger.Warn("Failed to save HTML file", &logging.LogFields{Error: err.Error()})
		} else {
			s.logger.Info("HTML file saved", &logging.LogFields{File: htmlFile})
		}
	}

	// Generate PDF
	if err := pdf.GeneratePDFChromedp(html, opts.OutputFile, s.logger); err != nil {
		s.logger.Error("PDF generation failed", &logging.LogFields{Error: err.Error()})
		return appErrs.NewPDFGenerationError("failed to generate PDF", err)
	}

	// ZUGFeRD compliance: generate and embed XML if requested
	if opts.EnableZUGFeRD {
		s.logger.Info("ZUGFeRD embedding requested", &logging.LogFields{File: opts.OutputFile})
		if opts.EmbeddedDataProvider != nil {
			s.logger.Info("Generating embedded data using provider", &logging.LogFields{File: opts.OutputFile})
			filePath, desc, err := opts.EmbeddedDataProvider.Generate(data, opts)
			if err != nil {
				s.logger.Error("Embedded data generation failed", &logging.LogFields{Error: err.Error(), File: opts.OutputFile})
				return appErrs.NewPDFGenerationError("failed to generate embedded data", err)
			}
			s.logger.Info("Embedded data generated", &logging.LogFields{File: filePath, Status: desc})
			err = compliance.EmbedZUGFeRDXML(opts.OutputFile, filePath, desc)
			if err != nil {
				s.logger.Error("Failed to embed data in PDF", &logging.LogFields{Error: err.Error(), File: opts.OutputFile})
				return appErrs.NewPDFGenerationError("failed to embed data in PDF", err)
			}
			s.logger.Info("Embedded data successfully added to PDF", &logging.LogFields{File: opts.OutputFile, Status: "ZUGFeRD XML embedded"})

			// 1. Validate ZUGFeRD XML against XSD
			xsdPath := "external/xrechnung_schema/resources/cii/16b/xsd/CrossIndustryInvoice_100pD16B.xsd" // Fixed path for XSD validation
			s.logger.Info("Validating embedded XML against XSD", &logging.LogFields{File: filePath, Status: xsdPath})
			xmlBytes, err := os.ReadFile(filePath)
			if err != nil {
				s.logger.Error("Failed to read embedded XML for validation", &logging.LogFields{Error: err.Error(), File: filePath})
				return appErrs.NewPDFGenerationError("failed to read embedded XML for validation", err)
			}
			err = xml.ValidateXMLWithSchema(xmlBytes, xsdPath)
			if err != nil {
				s.logger.Error("Embedded XML failed XSD validation", &logging.LogFields{Error: err.Error(), File: filePath})
				return appErrs.NewPDFGenerationError("embedded XML failed XSD validation", err)
			}
			s.logger.Info("Embedded XML passed XSD validation", &logging.LogFields{File: filePath, Status: "XSD validation passed"})

			// 2. Validate PDF/A-3u compliance using pdfcpu CLI
			s.logger.Info("Validating PDF/A-3u compliance using pdfcpu", &logging.LogFields{File: opts.OutputFile})
			pdfcpuCmd := exec.Command("pdfcpu", "validate", opts.OutputFile)
			output, err := pdfcpuCmd.CombinedOutput()
			if err != nil {
				s.logger.Error("PDF/A-3u compliance validation failed", &logging.LogFields{Error: string(output), File: opts.OutputFile})
				return appErrs.NewPDFGenerationError("PDF/A-3u compliance validation failed", fmt.Errorf("%s", string(output)))
			}
			if len(output) > 0 {
				s.logger.Info("PDF/A-3u validation output", &logging.LogFields{File: opts.OutputFile, Status: string(output)})
			}
		}
	}

	s.logger.Info("Invoice generated successfully", &logging.LogFields{File: opts.OutputFile})
	return nil
}

// ValidateInvoiceData validates invoice data without generating
func (s *InvoiceService) ValidateInvoiceData(data *models.InvoiceData) error {
	return s.validator.ValidateInvoiceData(data)
}

// GenerateInvoiceNumber creates a new invoice number based on the configured strategy
func (s *InvoiceService) GenerateInvoiceNumber() string {
	now := time.Now()
	
	switch s.config.Invoice.NumberingStrategy {
	case "timestamp":
		return fmt.Sprintf("%s%d", s.config.Invoice.NumberPrefix, now.Unix())
	case "date":
		return fmt.Sprintf("%s%s", s.config.Invoice.NumberPrefix, now.Format("2006-01-02"))
	case "sequential":
		// TODO: Implement sequential numbering with persistent storage
		return fmt.Sprintf("%s%s-001", s.config.Invoice.NumberPrefix, now.Format("2006-01"))
	default:
		return fmt.Sprintf("%s%s-001", s.config.Invoice.NumberPrefix, now.Format("2006-01"))
	}
}

// CreateSampleInvoice generates sample invoice data for testing
func (s *InvoiceService) CreateSampleInvoice() *models.InvoiceData {
	now := time.Now()
	dueDate := now.AddDate(0, 0, s.config.Invoice.DefaultDueDays)

	// Debug: Check config values
	s.logger.Debug("CreateSampleInvoice config values", &logging.LogFields{
		Currency: s.config.Invoice.DefaultCurrency,
		Lines: s.config.Invoice.DefaultDueDays,
	})

	currencySymbol := getCurrencySymbol(s.config.Invoice.DefaultCurrency)
	s.logger.Debug("Currency symbol lookup", &logging.LogFields{
		Currency: s.config.Invoice.DefaultCurrency,
		Status: currencySymbol,
	})
	
	invoiceData := &models.InvoiceData{
		Provider: models.CompanyInfo{
			ID:      uuid.New(),
			Name:    "Acme Corp Ltd",
			Address: models.Address{
				Street:     "123 Business Street, Suite 100",
				City:       "Business City",
				PostalCode: "12345",
				State:      "BC",
				Country:    "Canada",
			},
			VATID:     "CA123456789",
			Email:     "billing@acmecorp.com",
			Phone:     "+1 (555) 123-4567",
			Website:   "https://acmecorp.com",
			IBAN:      "CA29 NWBK 6016 1331 9268 19",
			SWIFT:     "NWBKCAXX",
		},
		Client: models.ClientInfo{
			ID:   uuid.New(),
			Name: "Tech Solutions Inc",
			Address: models.Address{
				Street:     "456 Client Avenue",
				City:       "Client City",
				PostalCode: "67890",
				State:      "CC",
				Country:    "USA",
			},
			Email: "accounts@techsolutions.com",
			Phone: "+1 (555) 987-6543",
			VATID: "US987654321",
		},
		Invoice: models.InvoiceDetails{
			ID:      uuid.New(),
			Number:  s.GenerateInvoiceNumber(),
			Date:    now,
			DueDate: dueDate,
			Status:  models.StatusDraft,
			Currency: models.Currency{
				Code:   s.config.Invoice.DefaultCurrency,
				Symbol: currencySymbol,
				Rate:   decimal.NewFromInt(1),
			},
			PaymentTerms: models.PaymentTerms{
				DueDays:     s.config.Invoice.DefaultDueDays,
				Description: "Payment terms: Net 30 days. Late payments subject to 1.5% monthly fee.",
			},
			Notes: "Professional services for Q3 2025",
			Lines: []models.InvoiceLine{
				{
					ID:          uuid.New(),
					Description: "Software Development Services",
					Quantity:    decimal.NewFromInt(40),
					UnitPrice:   decimal.NewFromFloat(125.00),
					TaxRate:     decimal.NewFromFloat(s.config.Invoice.DefaultTaxRate),
				},
				{
					ID:          uuid.New(),
					Description: "Technical Consultation",
					Quantity:    decimal.NewFromInt(8),
					UnitPrice:   decimal.NewFromFloat(150.00),
					TaxRate:     decimal.NewFromFloat(s.config.Invoice.DefaultTaxRate),
				},
				{
					ID:          uuid.New(),
					Description: "Project Management",
					Quantity:    decimal.NewFromInt(16),
					UnitPrice:   decimal.NewFromFloat(100.00),
					TaxRate:     decimal.NewFromFloat(s.config.Invoice.DefaultTaxRate),
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	
	// Calculate totals for the sample invoice
	invoiceData.Invoice.CalculateTotals()
	// Enable ZUGFeRD embedding for sample invoice
	invoiceData.EmbeddedData = models.EmbeddedDataZUGFeRD
	return invoiceData
}

// applyDefaults applies service configuration defaults to invoice data and options
func (s *InvoiceService) applyDefaults(data *models.InvoiceData, opts *GenerateOptions) {
	// Apply currency default - only if currency is not already set
	if opts.Currency != "" {
		// Override with command-line specified currency
		data.Invoice.Currency.Code = opts.Currency
		data.Invoice.Currency.Symbol = getCurrencySymbol(opts.Currency)
		if data.Invoice.Currency.Rate.IsZero() {
			data.Invoice.Currency.Rate = decimal.NewFromInt(1)
		}
	} else if data.Invoice.Currency.Code == "" {
		// Set default currency if none is specified
		data.Invoice.Currency = models.Currency{
			Code:   s.config.Invoice.DefaultCurrency,
			Symbol: getCurrencySymbol(s.config.Invoice.DefaultCurrency),
			Rate:   decimal.NewFromInt(1),
		}
	}
	// If currency is already set (like in sample invoice), leave it as is

	// Apply template default
	if opts.Template == "" || opts.Template == s.config.Template.Theme {
		opts.Template = ""
	}

	// Generate invoice number if not provided
	if data.Invoice.Number == "" {
		data.Invoice.Number = s.GenerateInvoiceNumber()
	}

	// Set due date if not provided
	if data.Invoice.DueDate.IsZero() {
		dueDate := time.Now().AddDate(0, 0, s.config.Invoice.DefaultDueDays)
		data.Invoice.DueDate = dueDate
	}

	// Set invoice date if not provided
	if data.Invoice.Date.IsZero() {
		data.Invoice.Date = time.Now()
	}

	// Apply default tax rate to lines that don't have one
	if opts.TaxRate != nil {
		defaultTaxRate := *opts.TaxRate
		for i := range data.Invoice.Lines {
			if data.Invoice.Lines[i].TaxRate.IsZero() {
				data.Invoice.Lines[i].TaxRate = decimal.NewFromFloat(defaultTaxRate)
			}
		}
	} else {
		for i := range data.Invoice.Lines {
			if data.Invoice.Lines[i].TaxRate.IsZero() {
				data.Invoice.Lines[i].TaxRate = decimal.NewFromFloat(s.config.Invoice.DefaultTaxRate)
			}
		}
	}

	// Apply custom tax rule if defined for invoice language/locale
	if rule := i18n.GetTaxRule(data.Invoice.Language); rule != nil {
		// Convert invoice data to map for rule
		invoiceMap := map[string]interface{}{
			"lines": data.Invoice.Lines,
			"subtotal": data.Invoice.Subtotal,
			"total_tax": data.Invoice.TotalTax,
			"currency": data.Invoice.Currency.Code,
			"country": data.Provider.Address.Country,
			// Add more fields as needed
		}
		customTax := rule(invoiceMap)
		data.Invoice.TotalTax = decimal.NewFromFloat(customTax)
		// TODO: Log and validate custom tax application
	}

	// In service layer, select provider based on data.EmbeddedData (from YAML)
	if data.EmbeddedData == models.EmbeddedDataZUGFeRD {
		opts.EmbeddedDataProvider = di.ProvidePDFEmbeddedDataProvider()
	} else {
		opts.EmbeddedDataProvider = nil
	}
}

// Interface for all ZUGFeRD XML builders
// Each profile should have its own builder implementing this
// Removed duplicate ZUGFeRDInvoiceXMLBuilder interface. Use from pkg/interfaces/interfaces.go

// getCurrencySymbol returns the symbol for a given currency code
func getCurrencySymbol(code string) string {
	symbols := map[string]string{
		"USD": "$", "EUR": "€", "GBP": "£", "JPY": "¥", "CHF": "CHF",
		"CAD": "C$", "AUD": "A$", "SEK": "kr", "NOK": "kr", "DKK": "kr",
		"PLN": "zł", "CZK": "Kč", "HUF": "Ft", "BGN": "лв", "RON": "lei",
		"HRK": "kn", "RUB": "₽", "CNY": "¥", "INR": "₹", "BRL": "R$",
		"MXN": "$", "ZAR": "R", "KRW": "₩", "SGD": "S$", "HKD": "HK$",
	}
	
	if symbol, exists := symbols[code]; exists {
		return symbol
	}
	
	// Default to currency code if symbol not found
	return code
}

// changeExtension changes the file extension while preserving the base name
func changeExtension(filename, newExt string) string {
	base := filename[:len(filename)-len(filepath.Ext(filename))]
	return base + newExt
}
