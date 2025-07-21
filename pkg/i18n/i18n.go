package i18n

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var locales = map[string]map[string]string{}

// Built-in tax rule types
var builtInTaxRuleTypes = map[string]func(params map[string]interface{}) TaxRule{
	"flat": func(params map[string]interface{}) TaxRule {
		flatRate, ok := params["rate"].(float64)
		if !ok {
			flatRate = 0.0
		}
		return func(data map[string]interface{}) float64 {
			if subtotal, ok := data["subtotal"].(float64); ok {
				return subtotal * flatRate
			}
			return 0.0
		}
	},
	// TODO: Add more built-in types (progressive, country-specific, etc.)
}

// LoadLocales loads the locale file for the given language
func LoadLocales(lang string) error {
	if _, ok := locales[lang]; ok {
		return nil
	}
	basePath := os.Getenv("LOCALES_PATH")
	if basePath == "" {
		basePath = "internal/i18n/locales"
	}
	path := filepath.Join(basePath, lang+".json")
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	var langMap map[string]interface{}
	dec := json.NewDecoder(f)
	if err := dec.Decode(&langMap); err != nil {
		return err
	}
	// Extract string translations
	strMap := map[string]string{}
	for k, v := range langMap {
		if s, ok := v.(string); ok {
			strMap[k] = s
		}
	}
	locales[lang] = strMap
	// Extract and register tax_rule if present
	if taxRuleRaw, ok := langMap["tax_rule"]; ok {
		if taxRuleMap, ok := taxRuleRaw.(map[string]interface{}); ok {
			typeName, _ := taxRuleMap["type"].(string)
			params, _ := taxRuleMap["params"].(map[string]interface{})
			if builder, ok := builtInTaxRuleTypes[typeName]; ok {
				RegisterTaxRule(lang, builder(params))
			}
		}
	}
	return nil
}

// GetTranslator returns a translation function for the given language
// with the ability to override specific keys
func GetTranslator(lang string, overrides map[string]string) func(string) string {
	_ = LoadLocales(lang)
	return func(key string) string {
		if overrides != nil {
			if v, ok := overrides[key]; ok {
				return v
			}
		}
		if m, ok := locales[lang]; ok {
			if v, ok := m[key]; ok {
				return v
			}
		}
		return key
	}
}

// TaxRule defines a custom tax calculation logic for a locale
// The function receives invoice data and returns the tax amount
// This allows per-country overrides via locale files or code
// TODO: Support loading logic from locale files (see TODO.md)
type TaxRule func(data map[string]interface{}) float64

// TaxRules registry for custom logic per locale
var TaxRules = map[string]TaxRule{}

// RegisterTaxRule allows setting a custom tax rule for a locale
func RegisterTaxRule(locale string, rule TaxRule) {
	TaxRules[locale] = rule
}

// GetTaxRule returns the tax rule for a locale, or nil if not set
func GetTaxRule(locale string) TaxRule {
	if rule, ok := TaxRules[locale]; ok {
		return rule
	}
	return nil
}
