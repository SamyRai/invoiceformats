// Package loader provides functions for loading and validating invoice data files.
package loader

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	appErrs "invoiceformats/pkg/errors"
	"invoiceformats/pkg/logging"
	"invoiceformats/pkg/models"
)

// LoadInvoiceData loads and validates invoice data from a YAML or JSON file
func LoadInvoiceData(filename string, logger logging.Logger) (*models.InvoiceData, error) {
	logger.Info("Attempting to load invoice file", &logging.LogFields{File: filename})
	cwd, cwdErr := os.Getwd()
	logger.Info("Current working directory", &logging.LogFields{Status: cwd, Error: func() string { if cwdErr != nil { return cwdErr.Error() } else { return "" } }()})
	if _, statErr := os.Stat(filename); statErr != nil {
		logger.Error("File does not exist or cannot be accessed", &logging.LogFields{File: filename, Error: statErr.Error()})
	} else {
		logger.Info("File exists and is accessible", &logging.LogFields{File: filename})
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Error("Failed to read invoice file", &logging.LogFields{File: filename, Error: err.Error()})
		return nil, appErrs.NewAppError(appErrs.ErrUnknown, "failed to read file", err)
	}

	var invoiceData models.InvoiceData
	ext := strings.ToLower(filepath.Ext(filename))
	var unmarshalErr error

	switch ext {
	case ".yaml", ".yml":
		unmarshalErr = yaml.Unmarshal(data, &invoiceData)
	case ".json":
		unmarshalErr = json.Unmarshal(data, &invoiceData)
	default:
		logger.Error("Unsupported file format", &logging.LogFields{File: filename, Error: "unsupported format", Status: ext})
		return nil, appErrs.NewAppError(appErrs.ErrUnknown, "unsupported file format (supported: .yaml, .yml, .json)", nil)
	}

	if unmarshalErr != nil {
		logger.Error("Failed to unmarshal invoice data", &logging.LogFields{File: filename, Error: unmarshalErr.Error()})
		return nil, appErrs.NewAppError(appErrs.ErrUnknown, "failed to parse invoice data", unmarshalErr)
	}

	// Log a summary of parsed data for analysis
	logger.Info("Parsed invoice data", &logging.LogFields{
		Provider:      invoiceData.Provider.Name,
		Client:        invoiceData.Client.Name,
		InvoiceNum:   invoiceData.Invoice.Number,
		Currency:      invoiceData.Invoice.Currency.Code,
		Lines:         len(invoiceData.Invoice.Lines),
		EmbeddedData: string(invoiceData.EmbeddedData),
	})

	// Input validation (basic)
	if invoiceData.Provider.Name == "" || invoiceData.Client.Name == "" || invoiceData.Invoice.Number == "" {
		logger.Error("Missing required fields in invoice data", &logging.LogFields{
			Provider:      invoiceData.Provider.Name,
			Client:        invoiceData.Client.Name,
			InvoiceNum:   invoiceData.Invoice.Number,
		})
		return nil, appErrs.NewAppError(appErrs.ErrValidationFailed, "missing required fields in invoice data", nil)
	}
	if len(invoiceData.Invoice.Lines) == 0 {
		logger.Error("Invoice data has no line items", &logging.LogFields{InvoiceNum: invoiceData.Invoice.Number})
		return nil, appErrs.NewAppError(appErrs.ErrValidationFailed, "invoice data must contain at least one line item", nil)
	}
	// TODO [high, 2h]: Add full struct validation using go-playground/validator for all fields
	// TODO [medium, 1h]: Validate custom types (decimal, uuid, time) for correct formats
	// TODO [low, 30m]: Add support for more file formats (e.g., TOML, XML)

	return &invoiceData, nil
}
