package loader

import (
	"os"
	"testing"

	"invoiceformats/testutils"

	"github.com/stretchr/testify/assert"
)

func writeTempFile(t *testing.T, content string, ext string) string {
	t.Helper()
	f, err := os.CreateTemp("", "invoice_test.*"+ext)
	assert.NoError(t, err)
	_, err = f.WriteString(content)
	assert.NoError(t, err)
	f.Close()
	return f.Name()
}

func TestLoadInvoiceData_ValidYAML(t *testing.T) {
	yaml := `
provider:
  name: "Test Provider"
  address:
    street: "123 Main St"
    city: "Testville"
    country: "US"
  email: "provider@example.com"
client:
  name: "Test Client"
  address:
    street: "456 Elm St"
    city: "Clienttown"
    country: "US"
  email: "client@example.com"
invoice:
  number: "INV-001"
  currency:
    code: "USD"
    symbol: "$"
    rate: 1.0
  lines:
    - description: "Service"
      quantity: 1
      unit_price: 100
      tax_rate: 10
      discount: 0
`
	file := writeTempFile(t, yaml, ".yaml")
	defer os.Remove(file)
	data, err := LoadInvoiceData(file, &testutils.TestLogger{})
	assert.NoError(t, err)
	assert.Equal(t, "Test Provider", data.Provider.Name)
	assert.Equal(t, "Test Client", data.Client.Name)
	assert.Equal(t, "INV-001", data.Invoice.Number)
	assert.Equal(t, "USD", data.Invoice.Currency.Code)
	assert.Len(t, data.Invoice.Lines, 1)
}

func TestLoadInvoiceData_MissingRequiredFields(t *testing.T) {
	yaml := `
provider:
  name: ""
client:
  name: ""
invoice:
  number: ""
  currency:
    code: "USD"
    symbol: "$"
    rate: 1.0
  lines: []
`
	file := writeTempFile(t, yaml, ".yaml")
	defer os.Remove(file)
	_, err := LoadInvoiceData(file, &testutils.TestLogger{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required fields")
}

func TestLoadInvoiceData_NoLines(t *testing.T) {
	yaml := `
provider:
  name: "Provider"
client:
  name: "Client"
invoice:
  number: "INV-002"
  currency:
    code: "USD"
    symbol: "$"
    rate: 1.0
  lines: []
`
	file := writeTempFile(t, yaml, ".yaml")
	defer os.Remove(file)
	_, err := LoadInvoiceData(file, &testutils.TestLogger{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must contain at least one line item")
}

func TestLoadInvoiceData_InvalidFormat(t *testing.T) {
	file := writeTempFile(t, "not valid", ".toml")
	defer os.Remove(file)
	_, err := LoadInvoiceData(file, &testutils.TestLogger{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported file format")
}

func TestLoadInvoiceData_MalformedYAML(t *testing.T) {
	yaml := `provider: [bad yaml`
	file := writeTempFile(t, yaml, ".yaml")
	defer os.Remove(file)
	_, err := LoadInvoiceData(file, &testutils.TestLogger{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse invoice data")
}

func TestLoadInvoiceData_MalformedJSON(t *testing.T) {
	json := `{ "provider": "bad json" `
	file := writeTempFile(t, json, ".json")
	defer os.Remove(file)
	_, err := LoadInvoiceData(file, &testutils.TestLogger{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse invoice data")
}
