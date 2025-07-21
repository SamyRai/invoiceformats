// Package zugferd provides ZUGFeRD-specific PDF embedding implementation.
package zugferd

import "invoiceformats/pkg/pdf"

// ZugferdEmbedder implements PDF embedding for ZUGFeRD XML.
type ZugferdEmbedder struct{}

func (e *ZugferdEmbedder) EmbedXML(pdf []byte, xml []byte, description string) ([]byte, error) {
	// TODO [context=zugferd pdf embedding, priority=high, effort=1h]: Implement PDF/A-3 embedding logic
	return nil, nil
}

var _ pdf.Embedder = (*ZugferdEmbedder)(nil)
