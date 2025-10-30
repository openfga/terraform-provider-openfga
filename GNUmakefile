default: fmt lint install generate

VERSION ?=
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
INSTALL_DIR := $(HOME)/.terraform.d/plugins/openfga/openfga/openfga/$(VERSION)/$(GOOS)_$(GOARCH)

# Standard Go build (for general development)
build:
	go build -v ./...

# Standard Go install (for general development)
install: build
	go install -v ./...

# Build provider binary for local Terraform testing
build-local:
ifndef VERSION
	$(error VERSION is required. Usage: make build-local VERSION=0.1.0)
endif
	mkdir -p bin
	go build -ldflags "-X main.version=$(VERSION)" -o bin/terraform-provider-openfga_v$(VERSION)

# Install provider for local Terraform testing
install-local: build-local
	mkdir -p $(INSTALL_DIR)
	cp bin/terraform-provider-openfga_v$(VERSION) $(INSTALL_DIR)/terraform-provider-openfga_v$(VERSION)

# Clean up local Terraform provider installations
clean:
	rm -rf $(HOME)/.terraform.d/plugins/openfga/openfga/openfga/

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m -p=1 ./...

.PHONY: fmt lint test testacc build install build-local install-local generate clean