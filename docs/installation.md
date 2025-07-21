# Installation

## Prerequisites

- Go 1.20+
- [pdfcpu](https://pdfcpu.io/) (for PDF/A-3 embedding)

## Steps

```sh
git clone https://github.com/your-org/invoiceformats.git
cd invoiceformats
go mod tidy
```

## Build CLI

```sh
go build -o invoicegen ./cmd/invoicegen
```
