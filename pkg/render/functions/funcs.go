package functions

import (
	"fmt"
	"html/template"
	"strings"
)

// NewTemplateFuncs returns a template.FuncMap with custom helpers and injected translator.
func NewTemplateFuncs(translator func(string) string) template.FuncMap {
	return template.FuncMap{
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
		"mod": func(a, b int) int { return a % b },
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
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
					return va > vb
				}
			case uint:
				if vb, ok := b.(uint); ok {
					return va > vb
				}
			case float64:
				if vb, ok := b.(float64); ok {
					return va > vb
				}
			case float32:
				if vb, ok := b.(float32); ok {
					return va > vb
				}
			}
			return false
		},
		"lt": func(a, b interface{}) bool {
			switch va := a.(type) {
			case int:
				if vb, ok := b.(int); ok {
					return va < vb
				}
			case uint:
				if vb, ok := b.(uint); ok {
					return va < vb
				}
			case float64:
				if vb, ok := b.(float64); ok {
					return va < vb
				}
			case float32:
				if vb, ok := b.(float32); ok {
					return va < vb
				}
			}
			return false
		},
		"eq": func(a, b interface{}) bool { return a == b },
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
		"t": translator,
	}
}
