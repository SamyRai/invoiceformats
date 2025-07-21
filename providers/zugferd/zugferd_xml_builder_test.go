package zugferd_test

import (
	"testing"

	"invoiceformats/pkg/models"
	"invoiceformats/providers/zugferd"
)

func mapToZUGFeRDInvoiceXML(inv models.ZUGFeRDInvoice) zugferd.ZUGFeRDInvoiceXML {
	return zugferd.ZUGFeRDInvoiceXML{
		XmlnsRsm: "urn:ferd:CrossIndustryDocument:invoice:1p0",
		XmlnsRam: "urn:un:unece:uncefact:data:standard:ReusableAggregateBusinessInformationEntity:12",
		XmlnsUdt: "urn:un:unece:uncefact:data:standard:UnqualifiedDataType:15",
		Context: zugferd.DocumentContextXML{GuidelineID: inv.Profile},
		Agreement: zugferd.TradeAgreementXML{
			Seller: zugferd.PartyXML{Name: inv.Seller.Name},
			Buyer:  zugferd.PartyXML{Name: inv.Buyer.Name},
		},
		Document: zugferd.DocumentXML{
			ID:        inv.DocumentID,
			IssueDate: zugferd.DateTimeXML{DateString: inv.IssueDate},
		},
		Transaction: zugferd.SupplyChainTradeTransactionXML{
			LineItems: []zugferd.LineItemXML{}, // TODO: map inv.LineItems
		},
		Settlement: zugferd.TradeSettlementXML{
			GrandTotal: inv.GrandTotal,
			Currency:   inv.Currency,
			Taxes:      []zugferd.TaxDetailXML{}, // TODO: map inv.Taxes
		},
	}
}

func TestBuildBasicXML_ValidInvoice(t *testing.T) {
	inv := models.ZUGFeRDInvoice{
		Profile:    "BASIC",
		Seller:     models.Party{Name: "Test Seller"},
		Buyer:      models.Party{Name: "Test Buyer"},
		DocumentID: "INV-123",
		IssueDate:  "20250713",
		GrandTotal: "100.00",
		Currency:   "EUR",
		LineItems: []models.LineItem{
			{
				Description: "Test Item",
				Quantity:    1,
				UnitPrice:   100.00,
				Total:       100.00,
			},
		},
		Taxes: []models.TaxDetail{
			{
				Type:   "VAT",
				Amount: 19.00,
				Rate:   19.0,
			},
		},
	}
	builder := zugferd.ZUGFeRDBasicXMLBuilder{}
	xmlInvoice := mapToZUGFeRDInvoiceXML(inv)
	xmlData, err := builder.BuildXML(xmlInvoice)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(xmlData) == 0 {
		t.Error("expected non-empty XML output")
	}
}

func TestBuildBasicXML_MissingFields(t *testing.T) {
	inv := models.ZUGFeRDInvoice{}
	builder := zugferd.ZUGFeRDBasicXMLBuilder{}
	xmlInvoice := mapToZUGFeRDInvoiceXML(inv)
	_, err := builder.BuildXML(xmlInvoice)
	if err == nil {
		t.Error("expected error for missing fields, got nil")
	}
}

func TestBuildBasicXML_InvalidDate(t *testing.T) {
	inv := models.ZUGFeRDInvoice{
		Seller:     models.Party{Name: "Seller"},
		Buyer:      models.Party{Name: "Buyer"},
		DocumentID: "INV-001",
		IssueDate:  "notadate",
		GrandTotal: "100.00",
		Currency:   "EUR",
	}
	builder := zugferd.ZUGFeRDBasicXMLBuilder{}
	xmlInvoice := mapToZUGFeRDInvoiceXML(inv)
	_, err := builder.BuildXML(xmlInvoice)
	if err == nil {
		t.Error("expected error for invalid IssueDate, got nil")
	}
}

// TODO: Add more tests for line items, taxes, and edge cases as model expands
