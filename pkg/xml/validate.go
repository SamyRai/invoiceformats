package xml

import (
	xsdvalidate "github.com/terminalstatic/go-xsd-validate"
)

// ValidateXMLWithSchema validates XML data against the provided XSD schema file.
// Returns a domain-specific error if validation fails.
func ValidateXMLWithSchema(xmlData []byte, xsdPath string) error {
	if err := xsdvalidate.Init(); err != nil {
		return &ValidationError{Reason: "failed to initialize XSD validator", Err: err}
	}
	defer xsdvalidate.Cleanup()

	xsdHandler, err := xsdvalidate.NewXsdHandlerUrl(xsdPath, xsdvalidate.ParsErrDefault)
	if err != nil {
		return &ValidationError{Reason: "failed to create XSD handler", Err: err}
	}
	defer xsdHandler.Free()

	if err := xsdHandler.ValidateMem(xmlData, xsdvalidate.ParsErrDefault); err != nil {
		return &ValidationError{Reason: "XML validation failed", Err: err}
	}
	return nil
}

// ValidationError represents an error during XML validation.
type ValidationError struct {
	Reason string
	Err    error
}

func (e *ValidationError) Error() string {
	return e.Reason + ": " + e.Err.Error()
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// TODO: [context: xml validation][priority: medium][effort: 1h] Add support for validating against multiple schemas if needed.
