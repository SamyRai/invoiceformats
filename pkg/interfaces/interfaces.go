// Package interfaces defines public interfaces for invoice generation and embedding.
package interfaces

// PDFEmbeddedDataProvider provides embedded data for PDF invoices.
type PDFEmbeddedDataProvider interface {
	// Generate creates embedded data for the invoice and returns the file path, description, and error.
	Generate(data interface{}, opts interface{}) (filePath string, description string, err error)
	// TODO: [context: interface, priority: low, effort: 15m] Consider splitting for different invoice types.
}

// ZUGFeRDInvoiceXMLBuilder builds ZUGFeRD-compliant XML for invoices.
type ZUGFeRDInvoiceXMLBuilder interface {
	BuildXML(invoiceData interface{}) ([]byte, error)
}

// TODO: [context: interface expansion, priority: low, effort: 30m] Add more interfaces for extensibility and future features.
