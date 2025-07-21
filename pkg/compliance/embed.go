// Package compliance provides domain-specific compliance logic for invoice formats.
package compliance

import (
	"bytes"
	"fmt"
	"invoiceformats/pkg/pdf"
	"io"
	"os"
)

// EmbedZUGFeRDXML embeds ZUGFeRD XML into a PDF file at the given path.
// Returns error if PDF or XML file is not found, or embedding fails.
func EmbedZUGFeRDXML(pdfPath, xmlPath, description string) error {
	pdfBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		return fmt.Errorf("failed to read PDF file: %w", err)
	}

	xmlBytes, err := os.ReadFile(xmlPath)
	if err != nil {
		return fmt.Errorf("failed to read XML file: %w", err)
	}

	// Use the PDF embedder for ZUGFeRD
	embedder := &pdf.ZugferdEmbedder{} // TODO [context=pdf embedding, priority=high, effort=2h]: Replace with ZUGFeRD-specific embedder if needed
	output, err := embedder.EmbedXML(pdfBytes, xmlBytes, description)
	if err != nil {
		return fmt.Errorf("failed to embed XML into PDF: %w", err)
	}

	// Overwrite the PDF file with the embedded result
	f, err := os.OpenFile(pdfPath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open PDF for writing: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(f, bytes.NewReader(output))
	if err != nil {
		return fmt.Errorf("failed to write embedded PDF: %w", err)
	}
	return nil
}

// Removed duplicate ValidateXMLWithSchema. Use from pkg/xml/validate.go
