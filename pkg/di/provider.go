package di

import (
	"fmt"
	"invoiceformats/pkg/interfaces"
	"invoiceformats/pkg/models"
	"invoiceformats/providers/zugferd"
	"os"
)

// ProvidePDFEmbeddedDataProvider returns a PDFEmbeddedDataProvider wired with DI
// This implementation directly uses the ZUGFeRD builder and domain mapping, no adapters.
type zugferdEmbeddedDataProvider struct {
	builder interfaces.ZUGFeRDInvoiceXMLBuilder
}

func (p *zugferdEmbeddedDataProvider) Generate(data models.InvoiceData, opts any) (string, string, error) {
	xmlBytes, err := p.builder.BuildXML(data)
	if err != nil {
		return "", "", fmt.Errorf("failed to build ZUGFeRD XML: %w", err)
	}
	// Save XML to temp file
	tmpFile, err := os.CreateTemp("", "zugferd-*.xml")
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()
	if _, err := tmpFile.Write(xmlBytes); err != nil {
		return "", "", fmt.Errorf("failed to write ZUGFeRD XML: %w", err)
	}
	return tmpFile.Name(), "ZUGFeRD XML Invoice", nil
}

func ProvidePDFEmbeddedDataProvider() interfaces.PDFEmbeddedDataProvider {
	builder := ProvideZUGFeRDInvoiceXMLBuilder()
	return &zugferdEmbeddedDataProvider{builder: builder}
}

// ProvideZUGFeRDInvoiceXMLBuilder returns the default ZUGFeRD XML builder implementation
func ProvideZUGFeRDInvoiceXMLBuilder() interfaces.ZUGFeRDInvoiceXMLBuilder {
	return zugferd.ZUGFeRDBasicXMLBuilder{}
}

// TODO [context: DI, priority: high, effort: medium]: Add more providers for other interfaces and implementations as refactor progresses
