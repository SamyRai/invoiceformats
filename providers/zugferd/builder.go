package zugferd

import (
	"encoding/xml"
	"errors"
	"fmt"
	"invoiceformats/pkg/models"
)

// InvoiceFormat specifies the XML output format.
type InvoiceFormat int

const (
	FormatZUGFeRD InvoiceFormat = iota
	FormatEN16931
)

// MapInvoiceToXML maps models.ZUGFeRDInvoice to ZUGFeRDInvoiceXML for XML generation.
// Supports both ZUGFeRD and EN16931 output formats.
func MapInvoiceToXML(inv models.ZUGFeRDInvoice, format InvoiceFormat) (ZUGFeRDInvoiceXML, error) {
	// Input validation
	if inv.DocumentID == "" {
		return ZUGFeRDInvoiceXML{}, errors.New("missing DocumentID")
	}
	if inv.IssueDate == "" {
		return ZUGFeRDInvoiceXML{}, errors.New("missing IssueDate")
	}
	if inv.Seller.Name == "" || inv.Buyer.Name == "" {
		return ZUGFeRDInvoiceXML{}, errors.New("missing Seller or Buyer name")
	}
	if inv.GrandTotal == "" || inv.Currency == "" {
		return ZUGFeRDInvoiceXML{}, errors.New("missing GrandTotal or Currency")
	}

	var rootName, nsRsm string
	if format == FormatEN16931 {
		rootName = "rsm:CrossIndustryInvoice"
		nsRsm = "urn:un:unece:uncefact:data:standard:CrossIndustryInvoice:100"
		// TODO [context: EN16931 compliance, priority: high, effort: 2h]: Ensure all EN16931 required fields and structure are present
	} else {
		rootName = "rsm:CrossIndustryInvoice"
		nsRsm = "urn:ferd:CrossIndustryDocument:invoice:1p0"
	}

	return ZUGFeRDInvoiceXML{
		XMLName:  xmlName(rootName),
		XmlnsRsm: nsRsm,
		XmlnsRam: "urn:un:unece:uncefact:data:standard:ReusableAggregateBusinessInformationEntity:12",
		XmlnsUdt: "urn:un:unece:uncefact:data:standard:UnqualifiedDataType:15",
		Context: DocumentContextXML{GuidelineID: inv.Profile},
		Agreement: TradeAgreementXML{
			Seller: mapParty(inv.Seller),
			Buyer:  mapParty(inv.Buyer),
		},
		Document: DocumentXML{
			ID:        inv.DocumentID,
			IssueDate: DateTimeXML{DateString: inv.IssueDate},
		},
		Transaction: SupplyChainTradeTransactionXML{
			LineItems: mapLineItems(inv.LineItems),
		},
		Settlement: TradeSettlementXML{
			GrandTotal: inv.GrandTotal,
			Currency:   inv.Currency,
			Taxes:      mapTaxDetails(inv.Taxes),
		},
	}, nil
}

// xmlName returns an xml.Name for the given local name.
func xmlName(local string) xml.Name {
	return xml.Name{Local: local}
}

func mapParty(p models.Party) PartyXML {
	return PartyXML{
		Name:    p.Name,
		VATID:   p.VATID, // Correctly map VATID
		Address: mapAddress(p.Address),
	}
}

func mapAddress(a models.Address) AddressXML {
	return AddressXML{
		Street:   a.Street,
		City:     a.City,
		PostCode: a.PostalCode, // Correctly map PostalCode to PostCode
		Country:  a.Country,
	}
}

func mapLineItems(items []models.LineItem) []LineItemXML {
	result := make([]LineItemXML, len(items))
	for i, item := range items {
		result[i] = LineItemXML{
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   fmt.Sprintf("%.2f", item.UnitPrice),
			Total:       fmt.Sprintf("%.2f", item.Total),
			TaxRate:     item.TaxRate, // Correctly map TaxRate
			// TODO [context: Line item XML, priority: medium, effort: medium]: Add product codes, units, etc.
		}
	}
	return result
}

func mapTaxDetails(taxes []models.TaxDetail) []TaxDetailXML {
	result := make([]TaxDetailXML, len(taxes))
	for i, tax := range taxes {
		result[i] = TaxDetailXML{
			Type:   tax.Type,
			Amount: fmt.Sprintf("%.2f", tax.Amount),
			Rate:   tax.Rate,
			// TODO [context: Tax details XML, priority: medium, effort: medium]: Add support for multi-rate VAT, exemptions, etc.
		}
	}
	return result
}

// TODO [context: builder.go, priority: high, effort: 1h]: Add unit tests for EN16931 output and input validation.
// TODO [context: builder.go, priority: medium, effort: 1h]: Document public API and edge cases.
