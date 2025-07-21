// Package xml provides interfaces and utilities for invoice XML generation and validation.
package xml

import "invoiceformats/pkg/models"

// Generator defines the interface for generating invoice XML.
type Generator interface {
	Generate(data models.InvoiceData) ([]byte, error)
}

// Validator defines the interface for validating invoice XML.
type Validator interface {
	Validate(xml []byte) error
}
