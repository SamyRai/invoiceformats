package interfaces

import (
	"invoiceformats/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRenderer struct{}

func (m *mockRenderer) Render(data models.InvoiceData, templatePath string) (string, error) {
	return "rendered", nil
}
func (m *mockRenderer) RenderOutputName(data models.InvoiceData, templatePath string) (string, error) {
	return "output_name", nil
}

type mockLocaleLoader struct{}

func (m *mockLocaleLoader) Load(lang, customPath string) (map[string]string, error) {
	return map[string]string{"hello": "world"}, nil
}

func TestTemplateRendererInterface(t *testing.T) {
	var r TemplateRenderer = &mockRenderer{}
	out, err := r.Render(models.InvoiceData{}, "")
	assert.NoError(t, err)
	assert.Equal(t, "rendered", out)
	name, err := r.RenderOutputName(models.InvoiceData{}, "")
	assert.NoError(t, err)
	assert.Equal(t, "output_name", name)
}

func TestLocaleLoaderInterface(t *testing.T) {
	var l LocaleLoader = &mockLocaleLoader{}
	loc, err := l.Load("en", "")
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"hello": "world"}, loc)
}
