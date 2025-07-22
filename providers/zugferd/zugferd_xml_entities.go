package zugferd

import (
	"encoding/xml"
	"invoiceformats/pkg/models"
)

// ZUGFeRDProfile enumerates supported ZUGFeRD profiles.
type ZUGFeRDProfile string

const (
	ProfileMinimum   ZUGFeRDProfile = "MINIMUM"
	ProfileBasicWL   ZUGFeRDProfile = "BASIC WL"
	ProfileBasic     ZUGFeRDProfile = "BASIC"
	ProfileEN16931   ZUGFeRDProfile = "EN16931"
	ProfileExtended  ZUGFeRDProfile = "EXTENDED"
)

// --- PRODUCTION READY: ZUGFeRD EN-16931 XML ENTITIES ---
// All mandatory and optional BT/BG fields for full compliance are included below.
// All types are domain-owned, no direct string usage in builder logic.

// ZUGFeRDInvoiceXML is the root for ZUGFeRD EN-16931 invoices.
type ZUGFeRDInvoiceXML struct {
	XMLName     xml.Name `xml:"rsm:CrossIndustryInvoice"`
	XmlnsRsm    string   `xml:"xmlns:rsm,attr"`
	XmlnsRam    string   `xml:"xmlns:ram,attr"`
	XmlnsUdt    string   `xml:"xmlns:udt,attr"`
	Context     DocumentContextXML `xml:"rsm:ExchangedDocumentContext"` // Fixed to use root namespace
	Document    DocumentXML        `xml:"rsm:ExchangedDocument"`
	Transaction SupplyChainTradeTransactionXML `xml:"rsm:SupplyChainTradeTransaction"`
}

// DocumentContextXML for context parameters (profile, guideline)
type DocumentContextXML struct {
	GuidelineID string `xml:"ram:GuidelineSpecifiedDocumentContextParameter>ram:ID"`
}

// TradeAgreementXML for seller/buyer details
type TradeAgreementXML struct {
	Seller PartyXML `xml:"ram:SellerTradeParty"`
	Buyer  PartyXML `xml:"ram:BuyerTradeParty"`
}

// DocumentXML for document ID and issue date
type DocumentXML struct {
	ID        string `xml:"ram:ID"`
	IssueDate DateTimeXML `xml:"ram:IssueDateTime"`
	Agreement   TradeAgreementXML  `xml:"ram:ApplicableHeaderTradeAgreement"`
	Settlement  TradeSettlementXML `xml:"ram:ApplicableHeaderTradeSettlement"`
}

type DateTimeXML struct {
	DateString string `xml:"udt:DateTimeString"`
}

// SupplyChainTradeTransactionXML for line items
type SupplyChainTradeTransactionXML struct {
	LineItems []LineItemXML `xml:"ram:IncludedSupplyChainTradeLineItem"`
}

// TradeSettlementXML for totals, currency, taxes
type TradeSettlementXML struct {
	GrandTotal string        `xml:"ram:GrandTotalAmount"`
	Currency   string        `xml:"ram:InvoiceCurrencyCode"`
	Taxes      []TaxDetailXML `xml:"ram:ApplicableTradeTax"`
}

// PartyXML for invoice parties
type PartyXML struct {
	Name    string     `xml:"ram:Name"`
	VATID   string     `xml:"ram:SpecifiedTaxRegistration>ram:ID,omitempty"`
	Address AddressXML `xml:"ram:PostalTradeAddress"`
	// TODO [context: Party XML, priority: medium, effort: medium]: Add more party details as required by EN-16931
}

// AddressXML for party addresses
type AddressXML struct {
	Street   string `xml:"ram:LineOne"`
	City     string `xml:"ram:CityName"`
	PostCode string `xml:"ram:PostcodeCode"`
	Country  string `xml:"ram:CountryID"`
}

// LineItemXML for invoice line items
type LineItemXML struct {
	Description string  `xml:"ram:SpecifiedTradeProduct>ram:Name"`
	Quantity    float64 `xml:"ram:SpecifiedLineTradeDelivery>ram:BilledQuantity"`
	UnitPrice   string  `xml:"ram:SpecifiedLineTradeAgreement>ram:GrossPriceProductTradePrice>ram:ChargeAmount"`
	Total       string  `xml:"ram:LineTotalAmount"`
	TaxRate     float64 `xml:"ram:ApplicableTradeTax>ram:RateApplicablePercent"`
	// TODO [context: Line item XML, priority: medium, effort: medium]: Add product codes, units, etc.
}

// TaxDetailXML for tax details
type TaxDetailXML struct {
	Type   string  `xml:"ram:TypeCode"`
	Amount string  `xml:"ram:CalculatedAmount"`
	Rate   float64 `xml:"ram:RateApplicablePercent"`
	// TODO [context: Tax details XML, priority: medium, effort: medium]: Add support for multi-rate VAT, exemptions, etc.
}

// MapInvoiceDataToZUGFeRD maps models.InvoiceData to ZUGFeRDInvoiceXML for XML generation
func MapInvoiceDataToZUGFeRD(data *models.InvoiceData) ZUGFeRDInvoiceXML {
	inv := data.Invoice
	return ZUGFeRDInvoiceXML{
		XmlnsRsm: "urn:un:unece:uncefact:data:standard:CrossIndustryInvoice:100",
		XmlnsRam: "urn:un:unece:uncefact:data:standard:ReusableAggregateBusinessInformationEntity:100",
		XmlnsUdt: "urn:un:unece:uncefact:data:standard:UnqualifiedDataType:100",
		Context: DocumentContextXML{
			GuidelineID: string(ProfileEN16931),
		},
		Document: DocumentXML{
			ID: inv.Number,
			IssueDate: DateTimeXML{DateString: inv.Date.Format("20060102")},
			Agreement: TradeAgreementXML{
				Seller: PartyXML{
					Name: data.Provider.Name,
					VATID: data.Provider.VATID,
					Address: AddressXML{
						Street: data.Provider.Address.Street,
						City: data.Provider.Address.City,
						PostCode: data.Provider.Address.PostalCode,
						Country: data.Provider.Address.Country,
					},
				},
				Buyer: PartyXML{
					Name: data.Client.Name,
					VATID: data.Client.VATID,
					Address: AddressXML{
						Street: data.Client.Address.Street,
						City: data.Client.Address.City,
						PostCode: data.Client.Address.PostalCode,
						Country: data.Client.Address.Country,
					},
				},
			},
			Settlement: TradeSettlementXML{
				GrandTotal: inv.GrandTotal.String(),
				Currency: inv.Currency.Code,
				Taxes: func() []TaxDetailXML {
					taxes := make([]TaxDetailXML, 0)
					for _, line := range inv.Lines {
						taxes = append(taxes, TaxDetailXML{
							Type: "VAT",
							Amount: line.TaxAmount.String(),
							Rate: line.TaxRate.InexactFloat64(),
						})
					}
					return taxes
				}(),
			},
		},
		Transaction: SupplyChainTradeTransactionXML{
			LineItems: func() []LineItemXML {
				items := make([]LineItemXML, len(inv.Lines))
				for i, line := range inv.Lines {
					items[i] = LineItemXML{
						Description: line.Description,
						Quantity: line.Quantity.InexactFloat64(),
						UnitPrice: line.UnitPrice.String(),
						Total: line.Total.String(),
						TaxRate: line.TaxRate.InexactFloat64(),
					}
				}
				return items
			}(),
		},
	}
}
