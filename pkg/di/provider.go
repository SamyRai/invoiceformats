package di

import (
	"fmt"
	"invoiceformats/pkg/interfaces"
	"invoiceformats/pkg/models"
	"invoiceformats/providers/zugferd"
	"os"
)

// ProvidePDFEmbeddedDataProvider returns a PDFEmbeddedDataProvider wired with DI
// This is a stub implementation. Replace with real provider as needed.
type basicEmbeddedDataProvider struct {
	builder interfaces.ZUGFeRDInvoiceXMLBuilder
}

func (p *basicEmbeddedDataProvider) Generate(data interface{}, opts interface{}) (string, string, error) {
	// TODO: [context: DI, priority: medium, effort: 1h] Implement real embedded data generation logic
	return "", "", nil
}

type zugferdEmbeddedDataProvider struct {
	builder interfaces.ZUGFeRDInvoiceXMLBuilder
}

func (p *zugferdEmbeddedDataProvider) Generate(data interface{}, opts interface{}) (string, string, error) {
	invoiceData, ok := data.(*models.InvoiceData)
	if !ok {
		return "", "", fmt.Errorf("invalid data type for ZUGFeRD provider: %T", data)
	}
	// Map InvoiceData to ZUGFeRDInvoiceXML (domain type)
	mapped := zugferd.MapInvoiceDataToZUGFeRD(invoiceData)
	xmlBytes, err := p.builder.BuildXML(mapped)
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

type zugferdBasicXMLBuilderAdapter struct {
	inner zugferd.ZUGFeRDBasicXMLBuilder
}

func (a zugferdBasicXMLBuilderAdapter) BuildXML(data interface{}) ([]byte, error) {
	invoice, ok := data.(zugferd.ZUGFeRDInvoiceXML)
	if !ok {
		return nil, fmt.Errorf("invalid type for ZUGFeRDInvoiceXMLBuilder: %T", data)
	}
	return a.inner.BuildXML(invoice)
}

// ProvideZUGFeRDInvoiceXMLBuilder returns the default ZUGFeRD XML builder implementation
func ProvideZUGFeRDInvoiceXMLBuilder() interfaces.ZUGFeRDInvoiceXMLBuilder {
	return zugferdBasicXMLBuilderAdapter{inner: zugferd.ZUGFeRDBasicXMLBuilder{}}
}

// TODO [context: DI, priority: high, effort: medium]: Add more providers for other interfaces and implementations as refactor progresses
