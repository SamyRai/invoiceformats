package xmlgen

import (
	"strings"
	"testing"

	"github.com/beevik/etree"
)

// ParseXML parses XML and returns the etree.Document, failing the test if invalid.
func ParseXML(t *testing.T, xmlData string) *etree.Document {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xmlData); err != nil {
		t.Fatalf("XML parse error: %v\nXML: %s", err, xmlData)
	}
	return doc
}

// AssertNamespace checks that the given namespace attribute is present and correct.
func AssertNamespace(t *testing.T, xmlData, ns, expected string) {
	if !strings.Contains(xmlData, ns+"=\""+expected+"\"") {
		t.Errorf("%s namespace missing or incorrect\nXML: %s", ns, xmlData)
	}
}

// FindElementByPath finds an element by local name path (ignoring prefixes).
func FindElementByPath(elem *etree.Element, path string) *etree.Element {
	parts := strings.Split(path, "/")
	current := elem
	for _, part := range parts {
		local := part
		if idx := strings.Index(local, ":"); idx != -1 {
			local = local[idx+1:]
		}
		found := false
		for _, child := range current.ChildElements() {
			if child.Tag == local {
				current = child
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	return current
}

// AssertElementValue checks that the element at the given path has the expected value (local names only).
func AssertElementValue(t *testing.T, doc *etree.Document, path, expected string) {
	elem := FindElementByPath(doc.Root(), path)
	if elem == nil {
		t.Errorf("Element %s not found", path)
		return
	}
	if elem.Text() != expected {
		t.Errorf("Element %s value mismatch: got %q, want %q", path, elem.Text(), expected)
	}
}

// AssertElementExists checks that the element at the given path exists (local names only).
func AssertElementExists(t *testing.T, doc *etree.Document, path string) {
	elem := FindElementByPath(doc.Root(), path)
	if elem == nil {
		t.Errorf("Element %s not found", path)
	}
}
