// Package validate provides command-line interface for the invoice generator.
package validate

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"invoiceformats/internal/config"
	appErrs "invoiceformats/pkg/errors"
	loader "invoiceformats/pkg/loader"
	"invoiceformats/pkg/logging"
	"invoiceformats/pkg/render"
	"invoiceformats/pkg/service"
)

var (
	validateFormat  string
	validateVerbose bool
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate [data-file]",
	Short: "Validate invoice data without generating PDF",
	Long: `Validate invoice data structure and business rules without generating a PDF.

This command checks:
• Required fields are present
• Data types are correct  
• Business rules are satisfied
• Currency codes are valid
• Email addresses are properly formatted
• VAT IDs follow correct patterns

Examples:
  # Validate YAML file
  invoicegen validate invoice-data.yaml

  # Validate JSON file with verbose output
  invoicegen validate invoice-data.json --verbose

  # Validate and show results in JSON format
  invoicegen validate data.yaml --format json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := GetLogger()
		inputFile := args[0]

		logger.Info("Validating invoice data", &logging.LogFields{File: inputFile})

		// Get invoice service
		invoiceService, err := GetInvoiceService()
		if err != nil {
			return fmt.Errorf("failed to create invoice service: %w", err)
		}

		// Load invoice data
		data, err := loader.LoadInvoiceData(inputFile, logger)
		if err != nil {
			return fmt.Errorf("failed to load invoice data: %w", err)
		}

		// Validate the invoice data
		if err := invoiceService.ValidateInvoiceData(data); err != nil {
			if appErr, ok := err.(*appErrs.AppError); ok {
				if validateFormat == "json" {
					fmt.Printf(`{"valid": false, "code": "%s", "error": "%s"}\n`, appErr.Code, appErr.Message)
				} else {
					fmt.Printf("❌ Validation failed [%s]: %s\n", appErr.Code, appErr.Message)
				}
				if validateVerbose && appErr.Cause != nil {
					fmt.Printf("Cause: %v\n", appErr.Cause)
				}
				return err
			}
			// fallback for non-AppError
			switch validateFormat {
			case "json":
				fmt.Printf(`{"valid": false, "error": "%s"}\n`, err.Error())
			default:
				fmt.Printf("❌ Validation failed: %s\n", err.Error())
			}
			return err
		}

		// Validation successful
		if validateVerbose {
			logger.Info("Validation successful", nil)
		}

		switch validateFormat {
		case "json":
			fmt.Printf(`{"valid": true, "message": "Invoice data is valid"}%s`, "\n")
		default:
			fmt.Printf("✅ Invoice data is valid\n")
			if validateVerbose {
				fmt.Printf("File: %s\n", inputFile)
				fmt.Printf("Provider: %s\n", data.Provider.Name)
				fmt.Printf("Client: %s\n", data.Client.Name)
				fmt.Printf("Invoice: %s\n", data.Invoice.Number)
				fmt.Printf("Lines: %d\n", len(data.Invoice.Lines))
				fmt.Printf("Currency: %s\n", data.Invoice.Currency.Code)
			}
		}

		return nil
	},
}

func init() {
	// Output format options
	validateCmd.Flags().StringVarP(&validateFormat, "format", "f", "text", "output format (text, json)")
	validateCmd.Flags().BoolVar(&validateVerbose, "verbose", false, "verbose validation output")
}

// ValidateCmd is the exported validate command
var ValidateCmd = validateCmd

// GetLogger returns a default logger instance
func GetLogger() logging.Logger {
	return logging.NewLogger()
}

// GetInvoiceService returns a default invoice service instance
func GetInvoiceService() (*service.InvoiceService, error) {
	cfg := config.DefaultConfig()
	logger := GetLogger()
	// Load embedded locales.json
	localeData, err := os.ReadFile("pkg/render/locales.json")
	if err != nil {
		return nil, fmt.Errorf("failed to load embedded locales.json: %w", err)
	}
	loader := render.NewLocaleLoader(localeData)
	return service.NewInvoiceService(cfg, logger, loader), nil
}
