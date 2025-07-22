package template

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"invoiceformats/pkg/models"
	"invoiceformats/pkg/render/interfaces"
	"path/filepath"
	"strings"
)

// Renderer implements TemplateRenderer for embedded and file-based templates.
type Renderer struct {
	TemplateFS    embed.FS
	TemplateFuncs template.FuncMap
	I18nProvider  interfaces.TranslatorProvider
}

func (r *Renderer) Render(data models.InvoiceData, templatePath string) (string, error) {
	lang := "en"
	if data.Invoice.Language != "" {
		lang = data.Invoice.Language
	}
	r.TemplateFuncs["t"] = r.I18nProvider(lang, nil)

	var tmpl *template.Template
	var err error
	if templatePath != "" {
		tmpl, err = template.New(filepath.Base(templatePath)).Funcs(r.TemplateFuncs).ParseFiles(templatePath)
		if err != nil {
			return "", err
		}
	} else {
		tmpl, err = template.New("invoice.html.tmpl").Funcs(r.TemplateFuncs).ParseFS(r.TemplateFS, "templates/invoice.html.tmpl")
		if err != nil {
			return "", err
		}
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (r *Renderer) RenderOutputName(data models.InvoiceData, templatePath string) (string, error) {
	lang := "en"
	if data.Invoice.Language != "" {
		lang = data.Invoice.Language
	}
	r.TemplateFuncs["t"] = r.I18nProvider(lang, nil)

	var tmpl *template.Template
	var err error
	if templatePath != "" {
		tmpl, err = template.New(filepath.Base(templatePath)).Funcs(r.TemplateFuncs).ParseFiles(templatePath)
		if err != nil {
			return "", err
		}
	} else {
		tmpl, err = template.New("invoice.html.tmpl").Funcs(r.TemplateFuncs).ParseFS(r.TemplateFS, "templates/invoice.html.tmpl")
		if err != nil {
			return "", err
		}
	}
	var buf bytes.Buffer
	if tmpl.Lookup("output_name") != nil {
		err = tmpl.ExecuteTemplate(&buf, "output_name", data)
		if err == nil {
			name := strings.TrimSpace(buf.String())
			if name != "" {
				return name, nil
			}
		}
	}
	return fmt.Sprintf("Invoice-%s", data.Invoice.Number), nil
}
