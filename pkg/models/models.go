package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Currency represents a currency with proper decimal handling
type Currency struct {
    Code   string          `json:"code" yaml:"code" validate:"required,len=3"`
    Symbol string          `json:"symbol" yaml:"symbol" validate:"required"`
    Rate   decimal.Decimal `json:"rate" yaml:"rate"`
}

// Address represents a structured address
type Address struct {
    Street     string `json:"street" yaml:"street" validate:"required"`
    City       string `json:"city" yaml:"city" validate:"required"`
    PostalCode string `json:"postal_code" yaml:"postal_code"`
    State      string `json:"state" yaml:"state"`
    Country    string `json:"country" yaml:"country" validate:"required"`
}

func (a Address) String() string {
    result := a.Street + "\n" + a.City
    if a.PostalCode != "" {
        result += " " + a.PostalCode
    }
    if a.State != "" {
        result += ", " + a.State
    }
    result += "\n" + a.Country
    return result
}

// CompanyInfo represents provider/company information
type CompanyInfo struct {
    ID          uuid.UUID `json:"id" yaml:"id"`
    Name        string    `json:"name" yaml:"name" validate:"required"`
    Address     Address   `json:"address" yaml:"address" validate:"required"`
    VATID       string    `json:"vat_id" yaml:"vat_id"`
    Email       string    `json:"email" yaml:"email" validate:"required,email"`
    Phone       string    `json:"phone" yaml:"phone"`
    Website     string    `json:"website" yaml:"website"`
    IBAN        string    `json:"iban" yaml:"iban"`
    SWIFT       string    `json:"swift" yaml:"swift"`
    TaxNumber   string    `json:"tax_number" yaml:"tax_number"`
    Logo        string    `json:"logo" yaml:"logo"` // Base64 or URL
}

// ClientInfo represents client/customer information  
type ClientInfo struct {
    ID      uuid.UUID `json:"id" yaml:"id"`
    Name    string    `json:"name" yaml:"name" validate:"required"`
    Address Address   `json:"address" yaml:"address" validate:"required"`
    Email   string    `json:"email" yaml:"email" validate:"required,email"`
    Phone   string    `json:"phone" yaml:"phone"`
    VATID   string    `json:"vat_id" yaml:"vat_id"`
}

// InvoiceLine represents a single line item on an invoice
type InvoiceLine struct {
    ID          uuid.UUID       `json:"id" yaml:"id"`
    Description string          `json:"description" yaml:"description" validate:"required"`
    Quantity    decimal.Decimal `json:"quantity" yaml:"quantity" validate:"required,gt=0"`
    UnitPrice   decimal.Decimal `json:"unit_price" yaml:"unit_price" validate:"required,gte=0"`
    Total       decimal.Decimal `json:"total" yaml:"total"`
    TaxRate     decimal.Decimal `json:"tax_rate" yaml:"tax_rate" validate:"gte=0,lte=100"`
    TaxAmount   decimal.Decimal `json:"tax_amount" yaml:"tax_amount"`
    Discount    decimal.Decimal `json:"discount" yaml:"discount" validate:"gte=0,lte=100"`
    Period      string          `json:"period" yaml:"period"` // For recurring/periodic services
}

// CalculateTotal calculates the total for this line item
func (il *InvoiceLine) CalculateTotal() {
    // Calculate base total
    baseTotal := il.Quantity.Mul(il.UnitPrice)
    
    // Apply discount if any
    if il.Discount.GreaterThan(decimal.Zero) {
        discountAmount := baseTotal.Mul(il.Discount).Div(decimal.NewFromInt(100))
        baseTotal = baseTotal.Sub(discountAmount)
    }
    
    // Calculate tax
    il.TaxAmount = baseTotal.Mul(il.TaxRate).Div(decimal.NewFromInt(100))
    
    // Set total
    il.Total = baseTotal.Add(il.TaxAmount)
}

// InvoiceStatus represents the status of an invoice
type InvoiceStatus string

const (
    StatusDraft     InvoiceStatus = "draft"
    StatusSent      InvoiceStatus = "sent"
    StatusPaid      InvoiceStatus = "paid"
    StatusOverdue   InvoiceStatus = "overdue"
    StatusCancelled InvoiceStatus = "cancelled"
)

// PaymentTerms represents payment terms
type PaymentTerms struct {
    DueDays     int    `json:"due_days" yaml:"due_days" validate:"gte=0"`
    Description string `json:"description" yaml:"description"`
}

// VATExemptionType represents types of VAT exemptions
type VATExemptionType string

const (
    VATExemptionNone        VATExemptionType = "none"
    VATExemptionSmallBusiness VATExemptionType = "small_business"
    VATExemptionNonEU       VATExemptionType = "non_eu"
    VATExemptionReverseCharge VATExemptionType = "reverse_charge"
    VATExemptionExport      VATExemptionType = "export"
    VATExemptionIntraCommunity VATExemptionType = "intra_community"
    VATExemptionEducation   VATExemptionType = "education"
    VATExemptionMedical     VATExemptionType = "medical"
    VATExemptionFinancial   VATExemptionType = "financial"
    VATExemptionOther       VATExemptionType = "other"
)

// TariffType represents types of tariffs
type TariffType string

const (
    TariffNone        TariffType = "none"
    TariffImport      TariffType = "import"
    TariffCustoms     TariffType = "customs"
    TariffExcise      TariffType = "excise"
    TariffEnvironmental TariffType = "environmental"
    TariffLuxury      TariffType = "luxury"
    TariffOther       TariffType = "other"
)

// EmbeddedDataType specifies what kind of data to embed in the PDF (e.g., "zugferd", "none").
type EmbeddedDataType string

const (
	EmbeddedDataNone    EmbeddedDataType = "none"
	EmbeddedDataZUGFeRD EmbeddedDataType = "zugferd"
	// TODO: Add more types as needed (e.g., xrechnung, peppol, custom)
)

// InvoiceDetails represents the main invoice information
type InvoiceDetails struct {
    ID           uuid.UUID       `json:"id" yaml:"id"`
    Number       string          `json:"number" yaml:"number" validate:"required"`
    Date         time.Time       `json:"date" yaml:"date"`
    DueDate      time.Time       `json:"due_date" yaml:"due_date"`
    Status       InvoiceStatus   `json:"status" yaml:"status"`
    Currency     Currency        `json:"currency" yaml:"currency" validate:"required"`
    Lines        []InvoiceLine   `json:"lines" yaml:"lines" validate:"required,min=1"`
    PaymentTerms PaymentTerms    `json:"payment_terms" yaml:"payment_terms"`
    Notes        string          `json:"notes" yaml:"notes"`
    Language     string          `json:"language" yaml:"language"` // Added for i18n
    Subtotal     decimal.Decimal `json:"subtotal" yaml:"subtotal"`
    TotalTax     decimal.Decimal `json:"total_tax" yaml:"total_tax"`
    TotalDiscount decimal.Decimal `json:"total_discount" yaml:"total_discount"`
    GrandTotal   decimal.Decimal `json:"grand_total" yaml:"grand_total"`
    CreatedAt    time.Time       `json:"created_at" yaml:"created_at"`
    UpdatedAt    time.Time       `json:"updated_at" yaml:"updated_at"`
    LegalFields  map[string]string `json:"legal_fields" yaml:"legal_fields"` // For country-specific legal requirements
    Recurrence   string          `json:"recurrence" yaml:"recurrence"` // Optional: invoice-level recurrence description
    VATExemptionType VATExemptionType `json:"vat_exemption_type" yaml:"vat_exemption_type"` // Use 'other' and set VATExemptionReason for custom
    VATExemptionReason string `json:"vat_exemption_reason" yaml:"vat_exemption_reason"` // Custom text if type is 'other'
    TariffType       TariffType       `json:"tariff_type" yaml:"tariff_type"` // Use 'other' and set AdditionalTariffs for custom
    AdditionalTariffs   string `json:"additional_tariffs" yaml:"additional_tariffs"`   // Custom text if type is 'other'
}

// CalculateTotals calculates all totals for the invoice
func (inv *InvoiceDetails) CalculateTotals() {
    var subtotal, totalTax, totalDiscount decimal.Decimal
    
    for i := range inv.Lines {
        // Generate ID if not set
        if inv.Lines[i].ID == uuid.Nil {
            inv.Lines[i].ID = uuid.New()
        }
        
        inv.Lines[i].CalculateTotal()
        
        lineSubtotal := inv.Lines[i].Quantity.Mul(inv.Lines[i].UnitPrice)
        
        // Calculate discount for this line
        if inv.Lines[i].Discount.GreaterThan(decimal.Zero) {
            lineDiscount := lineSubtotal.Mul(inv.Lines[i].Discount).Div(decimal.NewFromInt(100))
            totalDiscount = totalDiscount.Add(lineDiscount)
            lineSubtotal = lineSubtotal.Sub(lineDiscount)
        }
        
        subtotal = subtotal.Add(lineSubtotal)
        totalTax = totalTax.Add(inv.Lines[i].TaxAmount)
    }
    
    inv.Subtotal = subtotal
    inv.TotalTax = totalTax
    inv.TotalDiscount = totalDiscount
    inv.GrandTotal = subtotal.Add(totalTax)
}

// InvoiceData represents the complete invoice data structure
// Add EmbeddedDataType to allow specifying what to embed
type InvoiceData struct {
	Provider CompanyInfo    `json:"provider" yaml:"provider" validate:"required"`
	Client   ClientInfo     `json:"client" yaml:"client" validate:"required"`
	Invoice  InvoiceDetails `json:"invoice" yaml:"invoice" validate:"required"`
	EmbeddedData EmbeddedDataType `json:"embedded_data,omitempty" yaml:"embedded_data,omitempty"`
}
