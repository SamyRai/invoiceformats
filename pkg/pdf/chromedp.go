// Package pdf provides utilities for generating PDF files from HTML using headless Chrome (chromedp).
package pdf

import (
	"context"
	"errors"
	"io"
	"os"
	"sync"
	"time"

	"invoiceformats/pkg/logging"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var pdfLock sync.Mutex

// GeneratePDFChromedp renders the given HTML string to a PDF file using headless Chrome via chromedp.
// Returns error if PDF generation fails or output file cannot be written.
func GeneratePDFChromedp(html, outputPath string, logger logging.Logger) error {
	pdfLock.Lock()
	defer pdfLock.Unlock()

	if html == "" {
		return errors.New("input HTML is empty")
	}
	if outputPath == "" {
		return errors.New("output path is empty")
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Create a temporary HTML file
	tmpFile, tmpErr := os.CreateTemp("", "invoice-*.html")
	if tmpErr != nil {
		return tmpErr
	}
	defer os.Remove(tmpFile.Name())
	if _, writeErr := io.WriteString(tmpFile, html); writeErr != nil {
		tmpFile.Close()
		return writeErr
	}
	tmpFile.Close()

	var pdfBuf []byte
	url := "file://" + tmpFile.Name()

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	logger.Debug("Navigating to HTML file", &logging.LogFields{URL: url})

	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, pdfErr := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if pdfErr != nil {
				logger.Error("PrintToPDF error", &logging.LogFields{Error: pdfErr.Error()})
				return pdfErr
			}
			pdfBuf = buf
			return nil
		}),
	)
	if err != nil {
		logger.Error("chromedp.Run error", &logging.LogFields{Error: err.Error()})
		return err
	}

	if len(pdfBuf) == 0 {
		logger.Warn("Generated PDF is empty", nil)
		return errors.New("generated PDF is empty")
	}

	if err := os.WriteFile(outputPath, pdfBuf, 0644); err != nil {
		logger.Error("WriteFile error", &logging.LogFields{Error: err.Error()})
		return err
	}
	logger.Info("PDF written", &logging.LogFields{File: outputPath})
	return nil
}
