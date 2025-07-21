package render

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"invoiceformats/pkg/i18n"
	"invoiceformats/pkg/models"
	"path/filepath"
	"strings"
)

//go:embed templates/*
var templateFS embed.FS

// Generic comparison helpers for templates
func gt[T ~int|~uint|~float32|~float64](a, b T) bool { return a > b }
func lt[T ~int|~uint|~float32|~float64](a, b T) bool { return a < b }
func eq[T comparable](a, b T) bool { return a == b }

// templateFuncs provides custom template functions
var templateFuncs = template.FuncMap{
	"title": func(v interface{}) string {
		switch s := v.(type) {
		case string:
			return strings.Title(s)
		case fmt.Stringer:
			return strings.Title(s.String())
		default:
			return strings.Title(fmt.Sprintf("%v", s))
		}
	},
	"upper": func(v interface{}) string {
		switch s := v.(type) {
		case string:
			return strings.ToUpper(s)
		case fmt.Stringer:
			return strings.ToUpper(s.String())
		default:
			return strings.ToUpper(fmt.Sprintf("%v", s))
		}
	},
	"lower": func(v interface{}) string {
		switch s := v.(type) {
		case string:
			return strings.ToLower(s)
		case fmt.Stringer:
			return strings.ToLower(s.String())
		default:
			return strings.ToLower(fmt.Sprintf("%v", s))
		}
	},
	"capitalize": func(v interface{}) string {
		var s string
		switch val := v.(type) {
		case string:
			s = val
		case fmt.Stringer:
			s = val.String()
		default:
			s = fmt.Sprintf("%v", val)
		}
		if s == "" {
			return s
		}
		return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	},
	"mod": func(a, b int) int {
		return a % b
	},
	"add": func(a, b int) int {
		return a + b
	},
	"sub": func(a, b int) int {
		return a - b
	},
	"mul": func(a, b int) int {
		return a * b
	},
	"div": func(a, b int) int {
		if b == 0 {
			return 0
		}
		return a / b
	},
	"gt": func(a, b interface{}) bool {
		switch va := a.(type) {
		case int:
			if vb, ok := b.(int); ok {
				return gt(va, vb)
			}
		case uint:
			if vb, ok := b.(uint); ok {
				return gt(va, vb)
			}
		case float64:
			if vb, ok := b.(float64); ok {
				return gt(va, vb)
			}
		case float32:
			if vb, ok := b.(float32); ok {
				return gt(va, vb)
			}
		}
		return false
	},
	"lt": func(a, b interface{}) bool {
		switch va := a.(type) {
		case int:
			if vb, ok := b.(int); ok {
				return lt(va, vb)
			}
		case uint:
			if vb, ok := b.(uint); ok {
				return lt(va, vb)
			}
		case float64:
			if vb, ok := b.(float64); ok {
				return lt(va, vb)
			}
		case float32:
			if vb, ok := b.(float32); ok {
				return lt(va, vb)
			}
		}
		return false
	},
	"eq": func(a, b interface{}) bool {
		return a == b // Go's == works for comparable types
	},
	"default": func(value interface{}, def string) string {
		switch v := value.(type) {
		case nil:
			return def
		case string:
			if v == "" {
				return def
			}
			return v
		default:
			return fmt.Sprintf("%v", value)
		}
	},
}

func RenderHTML(data models.InvoiceData, templatePath string) (string, error) {
	lang := "en"
	if data.Invoice.Language != "" {
		lang = data.Invoice.Language
	}
	// Set up i18n translation function for templates
	templateFuncs["t"] = i18n.GetTranslator(lang, nil)

	var tmpl *template.Template
	var err error
	if templatePath != "" {
		tmpl, err = template.New(filepath.Base(templatePath)).Funcs(templateFuncs).ParseFiles(templatePath)
		if err != nil {
			return "", err
		}
	} else {
		// Use embedded template for default
		tmpl, err = template.New("invoice.html.tmpl").Funcs(templateFuncs).ParseFS(templateFS, "templates/invoice.html.tmpl")
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
