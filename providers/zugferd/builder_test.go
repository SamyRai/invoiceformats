package zugferd

import (
	"invoiceformats/pkg/models"
	"testing"
)

func TestMapInvoiceToXML_EN16931_Valid(t *testing.T) {
	inv := models.ZUGFeRDInvoice{
		Profile:    "EN16931",
		DocumentID: "INV-001",
		IssueDate:  "20250715",
		GrandTotal: "100.00",
		Currency:   "EUR",
		Seller: models.Party{
			Name:    "Test Seller",
			VATID:   "DE123456789",
			Address: models.Address{
				Street:     "Teststr. 1",
				City:       "Berlin",
				PostalCode: "10115",
				Country:    "DE",
			},
		},
		Buyer: models.Party{
			Name:    "Test Buyer",
			VATID:   "DE987654321",
			Address: models.Address{
				Street:     "Kaufstr. 2",
				City:       "Munich",
				PostalCode: "80331",
				Country:    "DE",
			},
		},
		LineItems: []models.LineItem{{
			Description: "Service",
			Quantity:    1,
			UnitPrice:   100.00,
			Total:       100.00,
			TaxRate:     19.0,
		}},
		Taxes: []models.TaxDetail{{
			Type:   "VAT",
			Amount: 19.00,
			Rate:   19.0,
		}},
	}

	xml, err := MapInvoiceToXML(inv, FormatEN16931)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if xml.XmlnsRsm != "urn:un:unece:uncefact:data:standard:CrossIndustryInvoice:100" {
		t.Errorf("expected EN16931 namespace, got %s", xml.XmlnsRsm)
	}
	if xml.Agreement.Seller.VATID != "DE123456789" {
		t.Errorf("expected Seller VATID, got %s", xml.Agreement.Seller.VATID)
	}
	if xml.Agreement.Seller.Address.PostCode != "10115" {
		t.Errorf("expected Seller PostCode, got %s", xml.Agreement.Seller.Address.PostCode)
	}
	if len(xml.Transaction.LineItems) == 0 || xml.Transaction.LineItems[0].TaxRate != 19.0 {
		t.Errorf("expected line item TaxRate 19.0, got %v", xml.Transaction.LineItems)
	}
}

func TestMapInvoiceToXML_MissingFields(t *testing.T) {
	inv := models.ZUGFeRDInvoice{}
	_, err := MapInvoiceToXML(inv, FormatEN16931)
	if err == nil {
		t.Error("expected error for missing fields, got nil")
	}
}

// TODO [context: builder_test.go, priority: medium, effort: 1h]: Add more tests for edge cases and ZUGFeRD output.
