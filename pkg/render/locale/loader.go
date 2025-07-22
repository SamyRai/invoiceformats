package locale

import (
	"encoding/json"
	"fmt"
	"os"
)

// Loader implements LocaleLoader for file-based and embedded sources.
type Loader struct {
	EmbeddedData []byte // If set, use this as the embedded locale data
}

// Load loads locales from embedded data or a custom file path.
func (l *Loader) Load(lang, customPath string) (map[string]string, error) {
	var locales map[string]map[string]string
	var data []byte
	var err error
	if customPath != "" {
		data, err = os.ReadFile(customPath)
		if err != nil {
			return nil, err
		}
	} else if l.EmbeddedData != nil {
		data = l.EmbeddedData
	} else {
		return nil, fmt.Errorf("no locale data provided")
	}
	if err := json.Unmarshal(data, &locales); err != nil {
		return nil, err
	}
	if l, ok := locales[lang]; ok {
		return l, nil
	}
	return nil, fmt.Errorf("locale %s not found", lang)
}
