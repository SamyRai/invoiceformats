package interfaces

import "invoiceformats/pkg/models"

// TemplateRenderer defines the interface for rendering templates.
type TemplateRenderer interface {
	Render(data models.InvoiceData, templatePath string) (string, error)
	RenderOutputName(data models.InvoiceData, templatePath string) (string, error)
}

// LocaleLoader defines the interface for loading locales.
type LocaleLoader interface {
	Load(lang, customPath string) (map[string]string, error)
}

// TranslatorProvider defines a function that returns a translation function for a given language and locales.
type TranslatorProvider func(lang string, locales map[string]string) func(string) string
