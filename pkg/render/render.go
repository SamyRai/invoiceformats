package render

import (
	"embed"
	"invoiceformats/pkg/models"
	"invoiceformats/pkg/render/functions"
	"invoiceformats/pkg/render/interfaces"
	"invoiceformats/pkg/render/locale"
	"invoiceformats/pkg/render/template"
)

//go:embed templates/* locales.json
var templateFS embed.FS

// NewRenderer creates a TemplateRenderer with embedded templates and injected i18n provider.
func NewRenderer(i18nProvider interfaces.TranslatorProvider) interfaces.TemplateRenderer {
	return &template.Renderer{
		TemplateFS:    templateFS,
		TemplateFuncs: functions.NewTemplateFuncs(i18nProvider("en", nil)), // default, will be replaced per call
		I18nProvider:  i18nProvider,
	}
}

// NewLocaleLoader creates a LocaleLoader with embedded locale data.
func NewLocaleLoader(embeddedData []byte) interfaces.LocaleLoader {
	return &locale.Loader{EmbeddedData: embeddedData}
}

// The following functions are for backward compatibility and simple use cases.
// For advanced use, use the interfaces directly.

// RenderHTML renders invoice data to HTML using embedded templates.
func RenderHTML(data models.InvoiceData, templatePath string, i18nProvider interfaces.TranslatorProvider) (string, error) {
	renderer := NewRenderer(i18nProvider)
	return renderer.Render(data, templatePath)
}

// RenderHTMLWithLocale renders invoice data to HTML with custom locale and template paths.
func RenderHTMLWithLocale(data models.InvoiceData, templatePath, localePath string, i18nProvider interfaces.TranslatorProvider, localeLoader interfaces.LocaleLoader) (string, error) {
	lang := data.Invoice.Language
	if lang == "" {
		lang = "en"
	}
	locales, err := localeLoader.Load(lang, localePath)
	if err != nil {
		return "", err
	}
	tFunc := func(key string) string {
		if v, ok := locales[key]; ok {
			return v
		}
		return key
	}
	renderer := &template.Renderer{
		TemplateFS:    templateFS,
		TemplateFuncs: functions.NewTemplateFuncs(tFunc),
		I18nProvider:  i18nProvider,
	}
	return renderer.Render(data, templatePath)
}

// RenderOutputName renders the output_name block or returns a default name.
func RenderOutputName(data models.InvoiceData, templatePath string, i18nProvider interfaces.TranslatorProvider) (string, error) {
	renderer := NewRenderer(i18nProvider)
	return renderer.RenderOutputName(data, templatePath)
}

// TODO [context: render.go, priority: medium, effort: 2h]: Refactor i18n package to support dependency injection for locale sources (embedded, file, remote)
// TODO [context: render.go, priority: low, effort: 1h]: Add support for multiple template types and dynamic selection
