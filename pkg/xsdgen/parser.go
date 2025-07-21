package xsdgen

// Parser is responsible for parsing XSD files into schema entities
// SRP: Only parses XSD, does not generate code
// Extensible via IParser interface

type IParser interface {
	Parse(path string) (*Structure, error)
}

type Parser struct{}

func (p *Parser) Parse(path string) (*Structure, error) {
	// TODO: Implement XSD parsing logic
	// Use encoding/xml to parse XSD and populate Structure/Entity/Field/Attribute
	return &Structure{}, nil
}

// TODO: Add unit tests for parser
