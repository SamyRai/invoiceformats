package locale

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoader_Load_Embedded(t *testing.T) {
	locales := map[string]map[string]string{
		"en": {"hello": "world"},
	}
	data, _ := json.Marshal(locales)
	loader := &Loader{EmbeddedData: data}
	res, err := loader.Load("en", "")
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"hello": "world"}, res)
}

func TestLoader_Load_CustomPath(t *testing.T) {
	// This test is a placeholder. To fully test, write a temp file and load it.
	// TODO [context: loader_test.go, priority: low, effort: 30m]: Add test for loading from customPath using a temp file
}

func TestLoader_Load_LocaleNotFound(t *testing.T) {
	locales := map[string]map[string]string{
		"en": {"hello": "world"},
	}
	data, _ := json.Marshal(locales)
	loader := &Loader{EmbeddedData: data}
	_, err := loader.Load("fr", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "locale fr not found")
}
