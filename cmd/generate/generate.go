// Package cmd provides command-line interface for the invoice generator.
package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"invoiceformats/internal/config"
	appErrs "invoiceformats/pkg/errors"
	"invoiceformats/pkg/loader"
	"invoiceformats/pkg/logging"
	"invoiceformats/pkg/models"
	"invoiceformats/pkg/service"
)

var (
	outputFile     string
	inputFile      string
	template       string
	currency       string
	taxRate        float64
	includeHTML    bool
	dryRun         bool
	validateOnly   bool
	sample         bool
)

// GetInvoiceService returns a default invoice service instance
func GetInvoiceService() (*service.InvoiceService, error) {
	cfg := config.DefaultConfig()
	logger := logging.NewLogger()
	return service.NewInvoiceService(cfg, logger), nil
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate [data-file]",
	Short: "Generate a PDF invoice",
	Long: `Generate a PDF invoice from the provided data.

The data can be provided as a YAML or JSON file, or use --sample to generate
a sample invoice with demo data.

Examples:
  # Generate from YAML file
  invoicegen generate invoice-data.yaml

  # Generate sample invoice
  invoicegen generate --sample

  # Generate with custom output file
  invoicegen generate --sample -o my-invoice.pdf

  # Validate data without generating PDF
  invoicegen generate data.yaml --validate-only

  # Dry run to see what would be generated
  invoicegen generate data.yaml --dry-run`,
	Args: func(cmd *cobra.Command, args []string) error {
		if !sample && len(args) == 0 {
			return fmt.Errorf("requires a data file argument or --sample flag")
		}
		if sample && len(args) > 0 {
			return fmt.Errorf("cannot specify data file when using --sample")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logging.NewLogger()

		// Get invoice service
		invoiceService, err := GetInvoiceService()
		if err != nil {
			return fmt.Errorf("failed to create invoice service: %w", err)
		}

		var data *models.InvoiceData

		if sample {
			logger.Info("Generating sample invoice", nil)
			data = invoiceService.CreateSampleInvoice()
			inputFile = "invoices/sample-invoice.yaml"
		} else {
			// Default to invoices/ dir if not absolute or already in invoices/
			if !filepath.IsAbs(args[0]) && !strings.HasPrefix(args[0], "invoices/") {
				inputFile = filepath.Join("invoices", args[0])
			} else {
				inputFile = args[0]
			}

			logger.Info("Loading invoice data", &logging.LogFields{File: inputFile})

			data, err = loader.LoadInvoiceData(inputFile, logger)
			if err != nil {
				return fmt.Errorf("failed to load invoice data: %w", err)
			}
		}

		// Set default output file if not specified
		if outputFile == "" {
			base := ""
			if sample {
				base = "sample-invoice"
			} else {
				base = strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
			}
			outputFile = filepath.Join("invoices/pdf", base+".pdf")
		}

		// Prepare generation options
		opts := &service.GenerateOptions{
			OutputFile:   outputFile,
			Template:     template,
			Currency:     currency,
			IncludeHTML:  includeHTML,
			DryRun:       dryRun,
			ValidateOnly: validateOnly,
		}

		if opts.Template == "" {
			// Use default embedded template if none specified
			opts.Template = "pkg/render/templates/invoice.html.tmpl"
		}

		if taxRate > 0 {
			opts.TaxRate = &taxRate
		}

		// Enable ZUGFeRD embedding if requested in YAML
		if data.EmbeddedData == models.EmbeddedDataZUGFeRD {
			opts.EnableZUGFeRD = true
		}

		// Generate invoice
		if err := invoiceService.GenerateInvoice(data, opts); err != nil {
			if appErr, ok := err.(*appErrs.AppError); ok {
				fmt.Fprintf(os.Stderr, "Error [%s]: %s\n", appErr.Code, appErr.Message)
				if appErr.Cause != nil {
					fmt.Fprintf(os.Stderr, "Cause: %v\n", appErr.Cause)
				}
				return err
			}
			return fmt.Errorf("failed to generate invoice: %w", err)
		}

		if validateOnly {
			logger.Info("Invoice data validation successful", nil)
		} else if dryRun {
			logger.Info("Dry run successful - would generate", &logging.LogFields{File: outputFile})
		} else {
			logger.Info("Invoice generated successfully", &logging.LogFields{File: outputFile})
		}

		return nil
	},
}

var GenerateCmd = generateCmd

func init() {
	// File options
	generateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output PDF file path")
	generateCmd.Flags().BoolVar(&sample, "sample", false, "generate sample invoice with demo data")

	// Template and styling options
	generateCmd.Flags().StringVarP(&template, "template", "t", "", "template theme to use")
	generateCmd.Flags().StringVarP(&currency, "currency", "c", "", "currency code (e.g., USD, EUR)")
	generateCmd.Flags().Float64Var(&taxRate, "tax-rate", 0, "default tax rate percentage")

	// Generation options
	generateCmd.Flags().BoolVar(&includeHTML, "include-html", false, "also save HTML output")
	generateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "validate and process but don't generate files")
	generateCmd.Flags().BoolVar(&validateOnly, "validate-only", false, "only validate the data, don't generate")
}
