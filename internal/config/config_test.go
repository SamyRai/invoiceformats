package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestAppConfig_Defaults(t *testing.T) {
	cfg := AppConfig{}
	assert.Zero(t, cfg.Invoice.DefaultCurrency)
	assert.Zero(t, cfg.Server.Port)
}

func TestInvoiceConfig_ValidationTags(t *testing.T) {
	cfg := InvoiceConfig{
		DefaultCurrency:   "USD",
		DefaultDueDays:    30,
		NumberingStrategy: "date-based",
		NumberPrefix:      "INV-",
		DefaultTaxRate:    10.0,
	}
	assert.Equal(t, "USD", cfg.DefaultCurrency)
	assert.Equal(t, 30, cfg.DefaultDueDays)
	assert.Equal(t, "date-based", cfg.NumberingStrategy)
	assert.Equal(t, 10.0, cfg.DefaultTaxRate)
}

func TestAppConfig_YAMLUnmarshal(t *testing.T) {
	yamlData := `
server:
  port: 8080
  host: "127.0.0.1"
  read_timeout: 10s
  write_timeout: 10s
  enable_cors: true
invoice:
  default_currency: "USD"
  default_due_days: 30
  numbering_strategy: "date"
  number_prefix: "INV-"
  default_tax_rate: 10.0
`
	var cfg AppConfig
	err := yaml.Unmarshal([]byte(yamlData), &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "USD", cfg.Invoice.DefaultCurrency)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "127.0.0.1", cfg.Server.Host)
	assert.Equal(t, 30, cfg.Invoice.DefaultDueDays)
	assert.Equal(t, "INV-", cfg.Invoice.NumberPrefix)
	assert.Equal(t, 10.0, cfg.Invoice.DefaultTaxRate)
	assert.Equal(t, "date", cfg.Invoice.NumberingStrategy)
	assert.Equal(t, 10*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, 10*time.Second, cfg.Server.WriteTimeout)
	assert.True(t, cfg.Server.EnableCORS)
}
