package zugferd

import (
	"encoding/xml"
	"fmt"

	xsdvalidate "github.com/terminalstatic/go-xsd-validate"
)

// ValidateXMLWithSchema validates XML data against the provided XSD schema file.
func ValidateXMLWithSchema(xmlData []byte, xsdPath string) error {
	if err := xsdvalidate.Init(); err != nil {
		return err
	}
	defer xsdvalidate.Cleanup()

	xsdHandler, err := xsdvalidate.NewXsdHandlerUrl(xsdPath, xsdvalidate.ParsErrDefault)
	if err != nil {
		return err
	}
	defer xsdHandler.Free()

	if err := xsdHandler.ValidateMem(xmlData, xsdvalidate.ParsErrDefault); err != nil {
		fmt.Println("\n--- XML Validation Failed ---")
		fmt.Println("XSD Path:", xsdPath)
		fmt.Println("Raw XML:")
		fmt.Println(string(xmlData))
		// Attempt pretty-print
		var prettyXML []byte
		var prettyErr error
		var v interface{}
		prettyErr = xml.Unmarshal(xmlData, &v)
		if prettyErr == nil {
			prettyXML, prettyErr = xml.MarshalIndent(v, "", "  ")
			if prettyErr == nil {
				fmt.Println("\nPretty-printed XML:")
				fmt.Println(string(prettyXML))
			}
		}
		fmt.Println("--- End XML Output ---")
		return err
	}
	return nil
}
