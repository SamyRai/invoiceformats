package render

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"invoiceformats/pkg/models"
	"invoiceformats/pkg/render/locale"
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

func fakeI18nProvider(lang string, locales map[string]string) func(string) string {
	return func(key string) string { return key }
}

func TestRenderHTML_Success(t *testing.T) {
	data := sampleInvoiceData()
	html, err := RenderHTML(data, "", fakeI18nProvider)
	require.NoError(t, err)
	assert.NotEmpty(t, html)
	assert.Contains(t, html, "Test Provider")
	assert.Contains(t, html, "Test Client")
	assert.Contains(t, html, "Service Fee")
	assert.Contains(t, html, "$")
}

func TestRenderHTML_InvalidData(t *testing.T) {
	// Provide incomplete data (missing required fields)
	data := models.InvoiceData{}
	html, err := RenderHTML(data, "", fakeI18nProvider)
	// The template will render zero values, but should not error
	assert.NoError(t, err)
	assert.NotEmpty(t, html)
	// Should not panic or return error, but output will be mostly empty
}

func TestRenderHTML_CustomTemplate(t *testing.T) {
	data := sampleInvoiceData()
	customPath := "test_custom.tmpl"
	tmplContent := `Custom Template: {{.Provider.Name}} - {{.Client.Name}}`
	os.WriteFile(customPath, []byte(tmplContent), 0644)
	defer os.Remove(customPath)
	html, err := RenderHTML(data, customPath, fakeI18nProvider)
	require.NoError(t, err)
	assert.Contains(t, html, "Custom Template")
	assert.Contains(t, html, data.Provider.Name)
	assert.Contains(t, html, data.Client.Name)
}

func TestRenderHTMLWithLocale_CustomTemplateAndLocale(t *testing.T) {
	data := sampleInvoiceData()
	customTemplate := "test_custom.tmpl"
	customLocale := "test_locales.json"
	tmplContent := `Custom Template: {{.Provider.Name}} - {{t "invoice"}}`
	localeContent := `{"en": {"invoice": "Test Invoice"}}`
	os.WriteFile(customTemplate, []byte(tmplContent), 0644)
	os.WriteFile(customLocale, []byte(localeContent), 0644)
	defer os.Remove(customTemplate)
	defer os.Remove(customLocale)

	// Use a real locale loader with embedded data
	var locales map[string]map[string]string
	_ = json.Unmarshal([]byte(localeContent), &locales)
	embeddedData, _ := json.Marshal(locales)
	loader := &locale.Loader{EmbeddedData: embeddedData}

	// Translation function that uses loaded locales
	translator := func(lang string, loc map[string]string) func(string) string {
		return func(key string) string {
			if v, ok := loc[key]; ok {
				return v
			}
			return key
		}
	}

	lang := data.Invoice.Language
	if lang == "" {
		lang = "en"
	}
	locMap, err := loader.Load(lang, customLocale)
	require.NoError(t, err)

	html, err := RenderHTMLWithLocale(data, customTemplate, customLocale, func(l string, _ map[string]string) func(string) string { return translator(l, locMap) }, loader)
	require.NoError(t, err)
	assert.Contains(t, html, "Custom Template")
	assert.Contains(t, html, data.Provider.Name)
	assert.Contains(t, html, "Test Invoice")
}

// TODO [context: render_test.go, priority: medium, effort: 1h]: Add more edge case tests for RenderHTMLWithLocale (invalid locale, missing template, etc.)
// Remove tests for embedded data type extraction, as compliance is now set from YAML, not template
