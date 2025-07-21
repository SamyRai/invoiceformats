package xsdgen

import (
	"invoiceformats/pkg/xsdgen/entities"
)

// XSDGen is the main struct for the code generator
// It follows SRP and SOLID principles: parsing, AST generation, and code output are separated.
type XSDGen struct {
	SchemaPath string
	Entities   []entities.Entity
}

// Structure holds the full schema structure for code generation
// SRP: Only responsible for holding schema structure

type Structure struct {
	Entities []entities.Entity
}

// TODO: Add parser, AST generator, and code output modules as separate files/classes
// TODO: Add interfaces for extensibility and dependency injection
