package xsdgen

import (
	"testing"
)

func TestParser_Parse(t *testing.T) {
	parser := &Parser{}
	structure, err := parser.Parse("testdata/sample.xsd")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if structure == nil {
		t.Error("Expected non-nil Structure")
	}
	// TODO: Add more assertions for parsed entities, fields, attributes
}
