package zugferd

import (
	"encoding/xml"
	"errors"
	"regexp"
)

// ZUGFeRDBasicXMLBuilder builds BASIC profile ZUGFeRD XML invoices using domain-owned types.
type ZUGFeRDBasicXMLBuilder struct{}

func (b ZUGFeRDBasicXMLBuilder) BuildXML(domain ZUGFeRDInvoiceXML) ([]byte, error) {
	// Input validation
	if domain.Context.GuidelineID == "" {
		return nil, errors.New("missing Profile (GuidelineID)")
	}
	if domain.Agreement.Seller.Name == "" {
		return nil, errors.New("missing Seller name")
	}
	if domain.Agreement.Buyer.Name == "" {
		return nil, errors.New("missing Buyer name")
	}
	if domain.Document.ID == "" {
		return nil, errors.New("missing DocumentID")
	}
	if domain.Document.IssueDate.DateString == "" {
		return nil, errors.New("missing IssueDate")
	}
	if !regexp.MustCompile(`^\d{8}$`).MatchString(domain.Document.IssueDate.DateString) {
		return nil, errors.New("invalid IssueDate format, expected YYYYMMDD")
	}
	if domain.Settlement.GrandTotal == "" {
		return nil, errors.New("missing GrandTotal")
	}
	if domain.Settlement.Currency == "" {
		return nil, errors.New("missing Currency")
	}
	return xml.MarshalIndent(domain, "", "  ")
}
