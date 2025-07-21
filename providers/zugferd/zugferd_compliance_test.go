package zugferd_test

import (
	"testing"

	"invoiceformats/pkg/models"
	"invoiceformats/pkg/xmlgen"
	"invoiceformats/providers/zugferd"
)

// Helper to generate XML for tests
func generateTestInvoiceXML() string {
	builder := zugferd.ZUGFeRDBasicXMLBuilder{}
	domain := models.ZUGFeRDInvoice{
		Profile:    "BASIC", // Ensure required Profile field is set
		Seller:     models.Party{Name: "Test Seller"},
		Buyer:      models.Party{Name: "Test Buyer"},
		DocumentID: "INV-001",
		IssueDate:  "20250714",
		GrandTotal: "100.00",
		Currency:   "EUR",
		LineItems:  []models.LineItem{{Description: "Service", Total: 100.00}},
	}
	// Map models.ZUGFeRDInvoice to zugferd.ZUGFeRDInvoiceXML before BuildXML
	xmlInvoice, err := zugferd.MapInvoiceToXML(domain, zugferd.FormatZUGFeRD)
	if err != nil {
		panic(err) // Fail test if XML generation fails
	}
	xmlBytes, err := builder.BuildXML(xmlInvoice)
	if err != nil {
		panic(err) // Fail test if XML generation fails
	}
	return string(xmlBytes)
}

func TestZUGFeRDRequiredAttributes(t *testing.T) {
	xmlData := generateTestInvoiceXML()
	doc := xmlgen.ParseXML(t, xmlData)
	root := doc.Root()
	if root == nil {
		t.Fatalf("Root element missing: want <rsm:CrossIndustryInvoice>\nXML: %s", xmlData)
	}
	if root.Tag != "CrossIndustryInvoice" {
		t.Errorf("Root tag incorrect: got '%s', want 'CrossIndustryInvoice'\nXML: %s", root.Tag, xmlData)
	}

	// Match actual builder output for namespace assertions
	xmlgen.AssertNamespace(t, xmlData, "xmlns:rsm", "urn:ferd:CrossIndustryDocument:invoice:1p0")
	xmlgen.AssertNamespace(t, xmlData, "xmlns:ram", "urn:un:unece:uncefact:data:standard:ReusableAggregateBusinessInformationEntity:12")
	xmlgen.AssertNamespace(t, xmlData, "xmlns:udt", "urn:un:unece:uncefact:data:standard:UnqualifiedDataType:15")

	// Use local names for element assertions
	xmlgen.AssertElementExists(t, doc, "ExchangedDocumentContext/GuidelineSpecifiedDocumentContextParameter/ID")
	xmlgen.AssertElementValue(t, doc, "ApplicableHeaderTradeSettlement/GrandTotalAmount", "100.00")
	// TODO: Add more business rule checks using helpers
}
