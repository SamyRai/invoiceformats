// Package zugferd provides the ZUGFeRD invoice provider implementation.
package zugferd

import (
	"invoiceformats/pkg/logging"
	"invoiceformats/pkg/models"
	"invoiceformats/pkg/pdf"
	xmlutil "invoiceformats/pkg/xml"
)

// ZUGFeRDProvider wires together XML builder, validator, and PDF embedder for ZUGFeRD invoices.
type ZUGFeRDProvider struct {
	Builder  ZUGFeRDBasicXMLBuilder
	Embedder pdf.Embedder
	Logger   logging.Logger
}

func NewZUGFeRDProvider(logger logging.Logger) *ZUGFeRDProvider {
	return &ZUGFeRDProvider{
		Builder:  ZUGFeRDBasicXMLBuilder{},
		Embedder: &ZugferdEmbedder{},
		Logger:   logger,
	}
}

// GenerateXML generates ZUGFeRD-compliant XML from invoice data.
func (p *ZUGFeRDProvider) GenerateXML(domain models.ZUGFeRDInvoice) ([]byte, error) {
	// Use new builder function for EN16931 compliance
	xmlInvoice, err := MapInvoiceToXML(domain, FormatEN16931)
	if err != nil {
		return nil, err
	}
	return p.Builder.BuildXML(xmlInvoice)
}

// ValidateXML validates ZUGFeRD XML against schemas and business rules.
func (p *ZUGFeRDProvider) ValidateXML(xmlData []byte, xsdPath string) error {
	return xmlutil.ValidateXMLWithSchema(xmlData, xsdPath)
}

// EmbedXMLIntoPDF embeds ZUGFeRD XML into a PDF document.
func (p *ZUGFeRDProvider) EmbedXMLIntoPDF(pdfBytes, xmlBytes []byte, description string) ([]byte, error) {
	return p.Embedder.EmbedXML(pdfBytes, xmlBytes, description)
}

// ...existing code for domain types should be moved to a separate file (e.g., builder.go) for single responsibility.
