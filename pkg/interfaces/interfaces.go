// Package interfaces defines public interfaces for invoice generation and embedding.
package interfaces

import "invoiceformats/pkg/models"

// PDFEmbeddedDataProvider provides embedded data for PDF invoices.
type PDFEmbeddedDataProvider interface {
	// Generate creates embedded data for the invoice and returns the file path, description, and error.
	Generate(data models.InvoiceData, opts any) (filePath string, description string, err error)
	// TODO: [context: interface, priority: low, effort: 15m] Consider splitting for different invoice types.
}

// ZUGFeRDInvoiceXMLBuilder builds ZUGFeRD-compliant XML for invoices.
type ZUGFeRDInvoiceXMLBuilder interface {
	BuildXML(invoiceData models.InvoiceData) ([]byte, error)
}

// TODO: [context: interface expansion, priority: low, effort: 30m] Add more interfaces for extensibility and future features.
