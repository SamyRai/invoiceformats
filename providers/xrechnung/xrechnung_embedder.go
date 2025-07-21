// Package xrechnung provides XRechnung-specific PDF embedding implementation.
package xrechnung

import "invoiceformats/pkg/pdf"

// XRechnungEmbedder implements PDF embedding for XRechnung XML.
type XRechnungEmbedder struct{}

func (e *XRechnungEmbedder) EmbedXML(pdf []byte, xml []byte, description string) ([]byte, error) {
	// TODO [context=xrechnung pdf embedding, priority=high, effort=1h]: Implement PDF/A-3 embedding logic
	return nil, nil
}

var _ pdf.Embedder = (*XRechnungEmbedder)(nil)
