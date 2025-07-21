// Package entities provides schema entity definitions for the XSD code generator.
package entities

// Entity represents a parsed XSD complex type or element
// SRP: One entity = one schema type
// SOLID: Extensible for new types

type Entity struct {
	Name       string
	Fields     []Field
	Attributes []Attribute
	Namespace  string
	IsComplex  bool
}

type Field struct {
	Name        string
	Type        string
	Namespace   string
	MinOccurs   int
	MaxOccurs   int
	IsAttribute bool
}

type Attribute struct {
	Name      string
	Type      string
	Namespace string
}
