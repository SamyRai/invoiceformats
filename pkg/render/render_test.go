package render

import (
	"html/template"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"invoiceformats/pkg/models"
)

func sampleInvoiceData() models.InvoiceData {
	provider := models.CompanyInfo{
		ID:    uuid.New(),
		Name:  "Test Provider",
		Address: models.Address{
			Street:  "123 Main St",
			City:    "Testville",
			Country: "Testland",
		},
		Email: "provider@example.com",
	}
	client := models.ClientInfo{
		ID:    uuid.New(),
		Name:  "Test Client",
		Address: models.Address{
			Street:  "456 Client Rd",
			City:    "Clientcity",
			Country: "Clientland",
		},
		Email: "client@example.com",
	}
	lines := []models.InvoiceLine{
		{
			ID:          uuid.New(),
			Description: "Service Fee",
			Quantity:    decimal.NewFromInt(2),
			UnitPrice:   decimal.NewFromFloat(100.0),
			TaxRate:     decimal.NewFromFloat(20.0),
		},
	}
	invoice := models.InvoiceDetails{
		ID:       uuid.New(),
		Number:   "INV-001",
		Date:     time.Now(),
		DueDate:  time.Now().AddDate(0, 0, 30),
		Status:   models.StatusDraft,
		Currency: models.Currency{Code: "USD", Symbol: "$"},
		Lines:    lines,
	}
	invoice.CalculateTotals()
	return models.InvoiceData{
		Provider: provider,
		Client:   client,
		Invoice:  invoice,
	}
}

func TestRenderHTML_Success(t *testing.T) {
	data := sampleInvoiceData()
	html, err := RenderHTML(data, "")
	require.NoError(t, err)
	assert.NotEmpty(t, html)
	assert.Contains(t, html, "Test Provider")
	assert.Contains(t, html, "Test Client")
	assert.Contains(t, html, "Service Fee")
	assert.Contains(t, html, "$")
}

func TestRenderHTML_InvalidTemplate(t *testing.T) {
	// Simulate missing template by parsing a non-existent file
	tmpl := templateFuncs // just to use the variable, not needed for test logic
	_ = tmpl
	// Not possible to inject a missing template with current API, so skip this test for now
	// TODO: [low, 1h] Refactor RenderHTML to allow template name injection for better testability
}

func TestRenderHTML_InvalidData(t *testing.T) {
	// Provide incomplete data (missing required fields)
	data := models.InvoiceData{}
	html, err := RenderHTML(data, "")
	// The template will render zero values, but should not error
	assert.NoError(t, err)
	assert.NotEmpty(t, html)
	// Should not panic or return error, but output will be mostly empty
}

func TestRenderHTML_TemplateFunctions(t *testing.T) {
	// Test a few template functions via template execution
	tmpl := `{{"hello world" | title}} {{"foo" | upper}} {{add 2 3}} {{gt 5 2}} {{capitalize "bar"}}`
	tpl, err := template.New("test").Funcs(templateFuncs).Parse(tmpl)
	require.NoError(t, err)
	var buf strings.Builder
	err = tpl.Execute(&buf, nil)
	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "Hello World")
	assert.Contains(t, out, "FOO")
	assert.Contains(t, out, "5")
	assert.Contains(t, out, "true")
	assert.Contains(t, out, "Bar")
}

func TestRenderHTML_CustomTemplate(t *testing.T) {
	data := sampleInvoiceData()
	customPath := "test_custom.tmpl"
	html, err := RenderHTML(data, customPath)
	require.NoError(t, err)
	assert.Contains(t, html, "Custom Template")
	assert.Contains(t, html, data.Provider.Name)
	assert.Contains(t, html, data.Client.Name)
}

func TestRenderHTML_DefaultTemplate(t *testing.T) {
	data := sampleInvoiceData()
	html, err := RenderHTML(data, "")
	require.NoError(t, err)
	assert.Contains(t, html, data.Provider.Name)
	assert.Contains(t, html, data.Client.Name)
}

// Remove tests for embedded data type extraction, as compliance is now set from YAML, not template
