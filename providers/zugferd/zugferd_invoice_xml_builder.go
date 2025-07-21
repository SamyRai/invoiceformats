package zugferd

// ZUGFeRDInvoiceXMLBuilder defines an interface for building ZUGFeRD XML.
type ZUGFeRDInvoiceXMLBuilder interface {
	BuildXML(domain ZUGFeRDInvoiceXML) ([]byte, error)
}
