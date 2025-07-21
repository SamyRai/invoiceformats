// Package invoiceformats provides abstract interfaces and types for invoice XML providers and compliance logic.
package invoiceformats

import "invoiceformats/pkg/models"

// InvoiceXMLProvider defines the interface for generating invoice XML for any format.
type InvoiceXMLProvider interface {
	// GenerateXML generates the XML representation of the invoice data.
	GenerateXML(data models.InvoiceData) ([]byte, error)
}

// ComplianceChecker defines the interface for validating invoice XML against format-specific rules.
type ComplianceChecker interface {
	// ValidateXML validates the XML and returns an error if non-compliant.
	ValidateXML(xml []byte) error
}

// PDFEmbedder defines the interface for embedding XML into a PDF document.
type PDFEmbedder interface {
	// EmbedXMLIntoPDF embeds the XML into the PDF and returns the result.
	EmbedXMLIntoPDF(pdf []byte, xml []byte, description string) ([]byte, error)
}

// Error types for domain-specific errors.
type ErrInvalidInvoice struct {
	Reason string
}

func (e ErrInvalidInvoice) Error() string {
	return "invalid invoice: " + e.Reason
}

// Dependency injection pattern: ProviderSet holds all dependencies for a format implementation.
type ProviderSet struct {
	XMLProvider       InvoiceXMLProvider
	ComplianceChecker ComplianceChecker
	PDFEmbedder       PDFEmbedder
}

// NewProviderSet constructs a ProviderSet for a given format.
func NewProviderSet(xml InvoiceXMLProvider, checker ComplianceChecker, embedder PDFEmbedder) *ProviderSet {
	return &ProviderSet{
		XMLProvider:       xml,
		ComplianceChecker: checker,
		PDFEmbedder:       embedder,
	}
}
