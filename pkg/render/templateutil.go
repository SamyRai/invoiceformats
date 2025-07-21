package render

import (
	"os"
	"path/filepath"
)

// TemplatePath returns the absolute path to the invoice template for robust loading
func TemplatePath() (string, error) {
	cwd, err := os.Getwd()
	var path string
	if err != nil {
		path = filepath.Join("pkg", "render", "templates", "invoice.html.tmpl") // fallback
	} else {
		path = filepath.Join(cwd, "pkg", "render", "templates", "invoice.html.tmpl")
	}
	if _, err := os.Stat(path); err != nil {
		return "", err // Return error if template does not exist
	}
	return path, nil
}

// TODO [context: template util, priority: medium, effort: low]: Consider supporting configurable template paths for multi-template scenarios
// TODO [context: test setup, priority: high, effort: low]: Ensure pkg/render/templates/invoice.html.tmpl is present in test environments. Add setup step or fixture if needed.
