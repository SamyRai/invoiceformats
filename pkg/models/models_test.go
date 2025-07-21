package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceLine_CalculateTotal(t *testing.T) {
	tests := []struct {
		name         string
		line         InvoiceLine
		expectedTotal decimal.Decimal
		expectedTax   decimal.Decimal
	}{
		{
			name: "basic calculation",
			line: InvoiceLine{
				Quantity:  decimal.NewFromInt(10),
				UnitPrice: decimal.NewFromFloat(100.0),
				TaxRate:   decimal.NewFromFloat(10.0),
				Discount:  decimal.Zero,
			},
			expectedTotal: decimal.NewFromFloat(1100.0), // 1000 + 100 tax
			expectedTax:   decimal.NewFromFloat(100.0),
		},
		{
			name: "with discount",
			line: InvoiceLine{
				Quantity:  decimal.NewFromInt(5),
				UnitPrice: decimal.NewFromFloat(200.0),
				TaxRate:   decimal.NewFromFloat(20.0),
				Discount:  decimal.NewFromFloat(10.0), // 10% discount
			},
			expectedTotal: decimal.NewFromFloat(1080.0), // 900 + 180 tax
			expectedTax:   decimal.NewFromFloat(180.0),
		},
		{
			name: "zero tax",
			line: InvoiceLine{
				Quantity:  decimal.NewFromInt(2),
				UnitPrice: decimal.NewFromFloat(50.0),
				TaxRate:   decimal.Zero,
				Discount:  decimal.Zero,
			},
			expectedTotal: decimal.NewFromFloat(100.0),
			expectedTax:   decimal.Zero,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.line.CalculateTotal()
			
			assert.True(t, tt.expectedTotal.Equal(tt.line.Total), 
				"Expected total %s, got %s", tt.expectedTotal, tt.line.Total)
			assert.True(t, tt.expectedTax.Equal(tt.line.TaxAmount), 
				"Expected tax %s, got %s", tt.expectedTax, tt.line.TaxAmount)
		})
	}
}

func TestInvoiceDetails_CalculateTotals(t *testing.T) {
	invoice := InvoiceDetails{
		ID:       uuid.New(),
		Number:   "TEST-001",
		Date:     time.Now(),
		DueDate:  time.Now().AddDate(0, 0, 30),
		Currency: Currency{Code: "USD", Symbol: "$", Rate: decimal.NewFromInt(1)},
		Lines: []InvoiceLine{
			{
				Description: "Item 1",
				Quantity:    decimal.NewFromInt(2),
				UnitPrice:   decimal.NewFromFloat(100.0),
				TaxRate:     decimal.NewFromFloat(10.0),
			},
			{
				Description: "Item 2",
				Quantity:    decimal.NewFromInt(1),
				UnitPrice:   decimal.NewFromFloat(50.0),
				TaxRate:     decimal.NewFromFloat(20.0),
				Discount:    decimal.NewFromFloat(10.0), // 10% discount
			},
		},
	}

	invoice.CalculateTotals()

	// Expected calculations:
	// Line 1: 2 * 100 = 200, tax = 20, total = 220
	// Line 2: 1 * 50 = 50, discount = 5, subtotal = 45, tax = 9, total = 54
	// Subtotal: 200 + 45 = 245
	// Total tax: 20 + 9 = 29
	// Grand total: 245 + 29 = 274
	// Total discount: 5

	expectedSubtotal := decimal.NewFromFloat(245.0)
	expectedTotalTax := decimal.NewFromFloat(29.0)
	expectedGrandTotal := decimal.NewFromFloat(274.0)
	expectedTotalDiscount := decimal.NewFromFloat(5.0)

	assert.True(t, expectedSubtotal.Equal(invoice.Subtotal), 
		"Expected subtotal %s, got %s", expectedSubtotal, invoice.Subtotal)
	assert.True(t, expectedTotalTax.Equal(invoice.TotalTax), 
		"Expected total tax %s, got %s", expectedTotalTax, invoice.TotalTax)
	assert.True(t, expectedGrandTotal.Equal(invoice.GrandTotal), 
		"Expected grand total %s, got %s", expectedGrandTotal, invoice.GrandTotal)
	assert.True(t, expectedTotalDiscount.Equal(invoice.TotalDiscount), 
		"Expected total discount %s, got %s", expectedTotalDiscount, invoice.TotalDiscount)
}

func TestAddress_String(t *testing.T) {
	tests := []struct {
		name     string
		address  Address
		expected string
	}{
		{
			name: "full address",
			address: Address{
				Street:     "123 Main St",
				City:       "New York",
				PostalCode: "10001",
				State:      "NY",
				Country:    "USA",
			},
			expected: "123 Main St\nNew York 10001, NY\nUSA",
		},
		{
			name: "minimal address",
			address: Address{
				Street:  "456 Oak Ave",
				City:    "Boston",
				Country: "USA",
			},
			expected: "456 Oak Ave\nBoston\nUSA",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.address.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCurrency_Validation(t *testing.T) {
	currency := Currency{
		Code:   "USD",
		Symbol: "$",
		Rate:   decimal.NewFromInt(1),
	}

	assert.Equal(t, "USD", currency.Code)
	assert.Equal(t, "$", currency.Symbol)
	assert.True(t, decimal.NewFromInt(1).Equal(currency.Rate))
}

func TestInvoiceData_Complete(t *testing.T) {
	// Create a complete invoice data structure
	now := time.Now()
	invoiceData := InvoiceData{
		Provider: CompanyInfo{
			ID:      uuid.New(),
			Name:    "Test Provider",
			Address: Address{
				Street:  "123 Provider St",
				City:    "Provider City",
				Country: "USA",
			},
			Email: "provider@test.com",
		},
		Client: ClientInfo{
			ID:   uuid.New(),
			Name: "Test Client",
			Address: Address{
				Street:  "456 Client Ave",
				City:    "Client City",
				Country: "USA",
			},
			Email: "client@test.com",
		},
		Invoice: InvoiceDetails{
			ID:       uuid.New(),
			Number:   "TEST-001",
			Date:     now,
			DueDate:  now.AddDate(0, 0, 30),
			Currency: Currency{Code: "USD", Symbol: "$", Rate: decimal.NewFromInt(1)},
			Lines: []InvoiceLine{
				{
					ID:          uuid.New(),
					Description: "Test Service",
					Quantity:    decimal.NewFromInt(1),
					UnitPrice:   decimal.NewFromFloat(100.0),
					TaxRate:     decimal.NewFromFloat(10.0),
				},
			},
		},
	}

	// Verify all fields are set
	require.NotEqual(t, uuid.Nil, invoiceData.Provider.ID)
	require.NotEqual(t, uuid.Nil, invoiceData.Client.ID)
	require.NotEqual(t, uuid.Nil, invoiceData.Invoice.ID)
	require.NotEmpty(t, invoiceData.Provider.Name)
	require.NotEmpty(t, invoiceData.Client.Name)
	require.NotEmpty(t, invoiceData.Invoice.Number)
	require.Len(t, invoiceData.Invoice.Lines, 1)

	// Calculate totals and verify
	invoiceData.Invoice.CalculateTotals()
	
	expectedTotal := decimal.NewFromFloat(110.0) // 100 + 10% tax
	assert.True(t, expectedTotal.Equal(invoiceData.Invoice.GrandTotal))
}
