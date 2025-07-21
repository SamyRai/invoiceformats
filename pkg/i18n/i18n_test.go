package i18n

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadLocalesAndTranslate(t *testing.T) {
	// Set LOCALES_PATH to absolute path for test reliability
	absPath, err := filepath.Abs("locales")
	if err != nil {
		t.Fatalf("failed to resolve absolute path: %v", err)
	}
	os.Setenv("LOCALES_PATH", absPath)

	langs := []string{"en", "de", "ru", "it"}
	for _, lang := range langs {
		err := LoadLocales(lang)
		if err != nil {
			t.Fatalf("failed to load locales for %s: %v", lang, err)
		}
		translator := GetTranslator(lang, nil)
		// Check a few keys
		if translator("invoice") == "invoice" {
			t.Errorf("translation for 'invoice' missing in %s", lang)
		}
		if translator("default_payment_terms") == "default_payment_terms" {
			t.Errorf("translation for 'default_payment_terms' missing in %s", lang)
		}
	}
}

func TestI18nOverride(t *testing.T) {
	absPath, err := filepath.Abs("locales")
	if err != nil {
		t.Fatalf("failed to resolve absolute path: %v", err)
	}
	os.Setenv("LOCALES_PATH", absPath)

	overrides := map[string]string{
		"invoice": "Custom Invoice",
		"default_payment_terms": "Custom payment terms!",
	}
	_ = LoadLocales("en")
	translator := GetTranslator("en", overrides)
	if got := translator("invoice"); got != "Custom Invoice" {
		t.Errorf("expected override for 'invoice', got %s", got)
	}
	if got := translator("default_payment_terms"); got != "Custom payment terms!" {
		t.Errorf("expected override for 'default_payment_terms', got %s", got)
	}
}
