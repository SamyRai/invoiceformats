package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTemplateFuncs_Basic(t *testing.T) {
	translator := func(key string) string { return "T:" + key }
	funcs := NewTemplateFuncs(translator)
	assert.NotNil(t, funcs["title"])
	assert.NotNil(t, funcs["t"])
	assert.Equal(t, "T:foo", funcs["t"].(func(string) string)("foo"))
}

func TestNewTemplateFuncs_Math(t *testing.T) {
	funcs := NewTemplateFuncs(func(key string) string { return key })
	add := funcs["add"].(func(int, int) int)
	assert.Equal(t, 5, add(2, 3))
	mod := funcs["mod"].(func(int, int) int)
	assert.Equal(t, 1, mod(7, 3))
}
