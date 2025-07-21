package models

// ZUGFeRDInvoice defines the structure for ZUGFeRD-compliant invoices.
type ZUGFeRDInvoice struct {
	Profile      string // MINIMUM, BASIC WL, BASIC, EN16931, EXTENDED
	Seller       Party
	Buyer        Party
	LineItems    []LineItem
	Taxes        []TaxDetail
	DocumentID   string
	IssueDate    string // YYYYMMDD
	GrandTotal   string
	Currency     string
	// TODO [context: ZUGFeRD profiles, priority: high, effort: medium]: Add support for all mandatory and optional EN-16931 fields
}

// Party represents an invoice party (seller or buyer).
type Party struct {
	Name    string
	VATID   string   // EN16931: VAT ID
	Address Address  // EN16931: Structured address (from models.go)
	// TODO [context: Party details, priority: medium, effort: low]: Add more party fields as required by EN-16931
}

// LineItem represents a single invoice line item.
type LineItem struct {
	Description string
	Quantity    float64
	UnitPrice   float64
	Total       float64
	TaxRate     float64 // EN16931: Tax rate for the line item
	// TODO [context: Line item details, priority: medium, effort: medium]: Add product codes, units, etc.
}

// TaxDetail represents tax information for the invoice.
type TaxDetail struct {
	Type   string
	Amount float64
	Rate   float64
	// TODO [context: Tax details, priority: medium, effort: low]: Add support for multi-rate VAT, exemptions, etc.
}

// ZUGFeRDInvoiceBuilder defines the interface for building ZUGFeRD invoices.
type ZUGFeRDInvoiceBuilder interface {
	BuildXML(inv ZUGFeRDInvoice) ([]byte, error)
}

// TODO [context: models, priority: high, effort: 1h]: Add unit tests for new fields and validation logic.
