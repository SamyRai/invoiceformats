package template

import (
	"embed"
	"invoiceformats/pkg/models"
	"invoiceformats/pkg/render/functions"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/invoice.html.tmpl
var testFS embed.FS

type fakeI18nProvider func(lang string, locales map[string]string) func(string) string

func (f fakeI18nProvider) Call(lang string, locales map[string]string) func(string) string {
	return f(lang, locales)
}

func TestRenderer_Render(t *testing.T) {
	r := &Renderer{
		TemplateFS:    testFS,
		TemplateFuncs: functions.NewTemplateFuncs(func(key string) string { return key }),
		I18nProvider:  func(lang string, locales map[string]string) func(string) string { return func(key string) string { return key } },
	}
	data := models.InvoiceData{Invoice: models.InvoiceDetails{Number: "INV-123"}}
	out, err := r.Render(data, "testdata/invoice.html.tmpl")
	assert.NoError(t, err)
	assert.Contains(t, out, "INV-123")
}

func TestRenderer_RenderOutputName(t *testing.T) {
	r := &Renderer{
		TemplateFS:    testFS,
		TemplateFuncs: functions.NewTemplateFuncs(func(key string) string { return key }),
		I18nProvider:  func(lang string, locales map[string]string) func(string) string { return func(key string) string { return key } },
	}
	data := models.InvoiceData{Invoice: models.InvoiceDetails{Number: "INV-456"}}
	name, err := r.RenderOutputName(data, "testdata/invoice.html.tmpl")
	assert.NoError(t, err)
	assert.Contains(t, name, "INV-456")
}
