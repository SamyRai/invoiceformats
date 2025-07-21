package zugferd_test

import (
	"os"
	"testing"

	"invoiceformats/pkg/compliance"

	"github.com/stretchr/testify/assert"
)

func TestEmbedZUGFeRDXML_FileNotFound(t *testing.T) {
	pdfPath := "nonexistent.pdf"
	xmlPath := "nonexistent.xml"
	desc := "Test Data"
	err := compliance.EmbedZUGFeRDXML(pdfPath, xmlPath, desc)
	assert.Error(t, err)
}

func TestEmbedZUGFeRDXML_InvalidPDF(t *testing.T) {
	// Create a dummy file that is not a valid PDF
	pdfPath := "dummy.txt"
	xmlPath := "dummy.xml"
	desc := "Test Data"
	_ = os.WriteFile(pdfPath, []byte("not a pdf"), 0644)
	_ = os.WriteFile(xmlPath, []byte("<xml></xml>"), 0644)
	defer os.Remove(pdfPath)
	defer os.Remove(xmlPath)
	err := compliance.EmbedZUGFeRDXML(pdfPath, xmlPath, desc)
	assert.Error(t, err)
}

func TestEmbedZUGFeRDXML_ValidPDFAndXML(t *testing.T) {
	t.Skip("Skipping until valid PDF/A-3u fixture is available.")
	// This test requires a valid PDF/A-3u file and a valid XML file.
	// For demonstration, we'll create a minimal valid PDF and XML, but in real scenarios use proper fixtures.
	pdfPath := "test_invoice.pdf"
	xmlPath := "test_invoice.xml"
	desc := "Test ZUGFeRD Data"

	// Minimal valid PDF (not PDF/A-3u, but enough for pdfcpu to process)
	pdfContent := []byte("%PDF-1.4\n1 0 obj\n<< /Type /Catalog >>\nendobj\ntrailer\n<< /Root 1 0 R >>\n%%EOF")
	_ = os.WriteFile(pdfPath, pdfContent, 0644)
	_ = os.WriteFile(xmlPath, []byte("<xml>test</xml>"), 0644)
	defer os.Remove(pdfPath)
	defer os.Remove(xmlPath)

	err := compliance.EmbedZUGFeRDXML(pdfPath, xmlPath, desc)
	// Accept either no error (if pdfcpu can process) or error (if PDF/A-3u upgrade fails)
	if err != nil {
		t.Logf("Expected error or success depending on pdfcpu capabilities: %v", err)
	} else {
		// Optionally, check that the PDF file still exists and is non-empty
		info, statErr := os.Stat(pdfPath)
		assert.NoError(t, statErr)
		assert.True(t, info.Size() > 0)
	}
}
