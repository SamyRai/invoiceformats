// Package schema provides interfaces and utilities for schema management and validation.
package schema

// Manager defines the interface for schema download, update, and validation.
type Manager interface {
	EnsureSchema(localPath string, remoteURL string) error
	Validate(xml []byte, schemaPath string) error
}
