# Architecture

- **cmd/**: CLI entrypoints
- **internal/**: config, schema, xmlgen utilities
- **pkg/**: core logic (models, pdf, render, validation, logging, i18n)
- **providers/**: format-specific logic (ZUGFeRD, XRechnung)
- **invoices/**: sample invoice data
- **external/**: schemas, XSLT, and resources

Follows SOLID, single responsibility, and clean architecture principles. Interfaces are defined in `pkg/interfaces/`.

See code comments and TODOs for further details.
