package pdf_test

import (
	"invoiceformats/pkg/logging"
	"invoiceformats/pkg/pdf"
	"invoiceformats/testutils"
	"os"
	"testing"
)

func getTestLogger() logging.Logger {
	return &testutils.TestLogger{}
}

func TestGeneratePDFChromedp_EmptyInput(t *testing.T) {
	err := pdf.GeneratePDFChromedp("", "out.pdf", getTestLogger())
	if err == nil {
		t.Error("expected error for empty HTML input")
	}
}

func TestGeneratePDFChromedp_EmptyOutputPath(t *testing.T) {
	err := pdf.GeneratePDFChromedp("<html></html>", "", getTestLogger())
	if err == nil {
		t.Error("expected error for empty output path")
	}
}

func TestGeneratePDFChromedp_BasicHTML(t *testing.T) {
	outPath := "test_invoice.pdf"
	err := pdf.GeneratePDFChromedp("<html><body><h1>Test</h1></body></html>", outPath, getTestLogger())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	info, statErr := os.Stat(outPath)
	if statErr != nil {
		t.Errorf("output file not created: %v", statErr)
	}
	if info.Size() == 0 {
		t.Error("output PDF is empty")
	}
	_ = os.Remove(outPath)
}
