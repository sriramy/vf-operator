##
## Make vf-operator, a network service that configures and discovers SR-IOV devices
##
## Targets;
##  help			This printout
##  all (default)		Build gRPC stubs and the executable
##  clean			Remove built files
##  dep			Installs pre-requisites
##  stubs			Generates gRPC stubs and OpenAPIv2 specs
##  test
##  install 		Installs the executable
##  swagger_install 	Installs static swagger UI)
##

BIN_DIR=bin
DOCS_DIR=docs
PROTO_DIR=pkg/api/v1
STUBS_DIR=$(PROTO_DIR)/gen
IMPORTS_DIR=$(PROTO_DIR)/imports
SWAGGER_DIR=swagger-ui
CMDS := $(patsubst ./cmd/%/,%,$(sort $(dir $(wildcard ./cmd/*/))))
BINARIES := $(patsubst %,bin/%,$(CMDS))

PREFIX ?= /usr/local
PROTOC ?= protoc
GO ?=go

all: dep stubs build

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

.PHONY: build
build: $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/ ./cmd/...

.PHONY: test
test: all
	$(GO) test ./...

.PHONY: check
check:
	$(GO) vet ./...

.PHONY: dep
dep:
	$(GO) install \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc \
		github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc \
		github.com/cweill/gotests/gotests

$(STUBS_DIR):
	mkdir -p $(STUBS_DIR)

$(DOCS_DIR):
	mkdir -p $(DOCS_DIR)

.PHONY: stubs
stubs: $(STUBS_DIR) $(DOCS_DIR)
	$(PROTOC) \
	-I $(PROTO_DIR) -I $(IMPORTS_DIR) \
	--go_out=$(STUBS_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(STUBS_DIR) --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=$(STUBS_DIR) --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=$(STUBS_DIR) \
	--doc_out=$(DOCS_DIR) --doc_opt=markdown,proto.md,source_relative \
	$(PROTO_DIR)/*/*.proto

.PHONY: clean
clean:
	$(GO) clean ./cmd/...
	rm -rf $(STUBS_DIR) 2>/dev/null
	rm -f $(BINARIES)

.PHONY: install
install:
	install -d $(DESTDIR)/$(PREFIX)/bin
	install -m 755 $(BIN_DIR)/vf-operator -t $(DESTDIR)/$(PREFIX)/bin

.PHONY: swagger_install
swagger_install:
	install -d $(DESTDIR)/$(PREFIX)/swagger-ui
	install -m 644 $(STUBS_DIR)/network/networkservice.swagger.json $(DESTDIR)/$(PREFIX)/swagger-ui/swagger.json
	install -m 644 $(SWAGGER_DIR)/* -t $(DESTDIR)/$(PREFIX)/swagger-ui

.PHONY: help
help:
	@grep '^##' $(lastword $(MAKEFILE_LIST)) | cut -c3-
	@echo "Binaries:"
	@echo "  $(BINARIES)"