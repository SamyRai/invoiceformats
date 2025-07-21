# InvoiceFormats

A modern, production-ready Go library and CLI for generating, validating, and embedding compliant electronic invoices (PDF, XML) in formats such as ZUGFeRD and XRechnung.

---

## Index

- [Overview](#overview)
- [Features](#features)
- [Supported Formats](#supported-formats)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [Configuration](#configuration)
- [Architecture](#architecture)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)
- [Roadmap](#roadmap)
- [Changelog](#changelog)

---

## Overview

InvoiceFormats provides a robust, extensible platform for generating, validating, and embedding electronic invoices in compliance with European standards. It supports both CLI and API usage, with a focus on modularity, testability, and clean architecture.

## Features

- Generate PDF invoices with embedded XML (ZUGFeRD, XRechnung)
- Validate invoices against EN16931 and UBL schemas
- Modern HTML templates with i18n support
- Dependency injection and interface-driven design
- Comprehensive error handling and logging
- CLI and HTTP API modes
- Docker containerization support

## Supported Formats

- **ZUGFeRD** (all profiles)
- **XRechnung**
- **EN16931**
- **UBL** (via XSLT)

## Getting Started

### Prerequisites

- Go 1.20+
- [pdfcpu](https://pdfcpu.io/) (for PDF/A-3 embedding)

### Installation

```sh
git clone https://github.com/your-org/invoiceformats.git
cd invoiceformats
go mod tidy
```

### Build CLI

```sh
go build -o invoicegen ./cmd/invoicegen
```

## Usage

### CLI

```sh
./invoicegen generate --input invoices/sample-invoice.yaml --output out.pdf
./invoicegen validate --input out.pdf
```

### Library

Import and use Go interfaces from `pkg/` and `providers/` for custom integration.

## Configuration

- See `internal/config/config.go` for YAML/ENV config structure.
- Example config options: default currency, due days, PDF settings, template theme, etc.

## Architecture

- **cmd/**: CLI entrypoints
- **internal/**: config, schema, xmlgen utilities
- **pkg/**: core logic (models, pdf, render, validation, logging, i18n)
- **providers/**: format-specific logic (ZUGFeRD, XRechnung)
- **invoices/**: sample invoice data
- **external/**: schemas, XSLT, and resources

Follows SOLID, single responsibility, and clean architecture principles. Interfaces are defined in `pkg/interfaces/`.

## Testing

- Run all tests: `go test ./...`
- Unit tests cover all major logic, including edge cases and error handling.
- See `providers/zugferd/builder_test.go`, `pkg/render/render_test.go`, etc.

## Contributing

- Please see [CONTRIBUTING.md](CONTRIBUTING.md) (to be created)
- Follow code style, add tests, and update TODOs as needed
- All contributions must include tests and documentation

## License

Specify your license in a `LICENSE` file (MIT, Apache 2.0, etc.)

## Roadmap

- [ ] Add HTTP API server
- [ ] Expand UBL/EN16931 validation
- [ ] Add more invoice templates and i18n locales
- [ ] Improve error domain and logging
- [ ] Add more sample invoices and docs

## Changelog

See `CHANGELOG.md` (to be created)

---

For more details, see code comments, TODOs, and the `pkg/` and `providers/` directories.
