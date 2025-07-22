# Makefile for invoiceformats

APP_NAME=invoicegen
CMD_PATH=cmd/invoicegen/main.go
TEMPLATE_SRC_DIRS=pkg/render/templates pkg/service
TEMPLATE_DST_DIR=invoices/templates

.PHONY: all lint build generate generate-templates

all: lint build generate generate-templates

lint:
	@if [ -f .golangci.yml ] || [ -f .golangci.yaml ]; then \
		golangci-lint run ./... ; \
	else \
		go vet ./... ; \
	fi

build:
	go build -o bin/$(APP_NAME) $(CMD_PATH)

generate:
	go generate ./...

# Copy all .tmpl files from template source dirs to a central location for inspection/testing
# (adjust as needed for your actual template generation process)
generate-templates:
	@mkdir -p $(TEMPLATE_DST_DIR)
	@find $(TEMPLATE_SRC_DIRS) -name '*.tmpl' -exec cp {} $(TEMPLATE_DST_DIR) \;
	@echo "Templates copied to $(TEMPLATE_DST_DIR)"
	@mkdir -p invoices/pdf
	@for tmpl in $(TEMPLATE_DST_DIR)/*.tmpl; do \
	  name=$$(basename $$tmpl .html.tmpl); \
	  echo "Generating PDF for template: $$tmpl"; \
	  ./bin/$(APP_NAME) generate invoices/sample-invoice.yaml --template "$$tmpl" --output "invoices/pdf/$$name.pdf" || echo "Failed to generate PDF for $$tmpl"; \
	done
	@echo "PDFs generated in invoices/pdf/"

# Generate a ZUGFeRD PDF using the glpx.yaml invoice and the standard invoice template (with shadow removed)
generate-zugferd:
	@mkdir -p invoices/pdf
	./bin/$(APP_NAME) generate invoices/glpx.yaml --template pkg/render/templates/invoice.html.tmpl --output invoices/pdf/glpx-zugferd.pdf || echo "Failed to generate ZUGFeRD PDF"
