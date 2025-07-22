package zugferd

import (
	"encoding/xml"
	"errors"
	"invoiceformats/pkg/models"
	"regexp"
)

// ZUGFeRDBasicXMLBuilder builds BASIC profile ZUGFeRD XML invoices using domain-owned types.
type ZUGFeRDBasicXMLBuilder struct{}

// BuildXML implements interfaces.ZUGFeRDInvoiceXMLBuilder and takes models.InvoiceData directly.
func (b ZUGFeRDBasicXMLBuilder) BuildXML(data models.InvoiceData) ([]byte, error) {
	// Check for zero date before mapping
	if data.Invoice.Date.IsZero() {
		return nil, errors.New("invalid IssueDate: zero value")
	}
	mapped := MapInvoiceDataToZUGFeRD(&data)
	// Input validation (copied from old BuildXML)
	if mapped.Context.GuidelineID == "" {
		return nil, errors.New("missing Profile (GuidelineID)")
	}
	if mapped.Agreement.Seller.Name == "" {
		return nil, errors.New("missing Seller name")
	}
	if mapped.Agreement.Buyer.Name == "" {
		return nil, errors.New("missing Buyer name")
	}
	if mapped.Document.ID == "" {
		return nil, errors.New("missing DocumentID")
	}
	if mapped.Document.IssueDate.DateString == "" {
		return nil, errors.New("missing IssueDate")
	}
	if !regexp.MustCompile(`^\d{8}$`).MatchString(mapped.Document.IssueDate.DateString) {
		return nil, errors.New("invalid IssueDate format, expected YYYYMMDD")
	}
	if mapped.Settlement.GrandTotal == "" {
		return nil, errors.New("missing GrandTotal")
	}
	if mapped.Settlement.Currency == "" {
		return nil, errors.New("missing Currency")
	}
	return xml.MarshalIndent(mapped, "", "  ")
}
