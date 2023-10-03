BIN_DIR=bin
PROTO_DIR=pkg/api
STUBS_DIR=pkg/stubs

PROTOC=protoc
GO=go

all: stubs build

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

.PHONY: build
build: $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/vf-operator ./cmd/...

.PHONY: dep
dep:
	$(GO) mod tidy
	$(GO) install \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc

$(STUBS_DIR):
	mkdir -p $(STUBS_DIR)

.PHONY: stubs
stubs: $(STUBS_DIR)
	$(PROTOC) \
	-I $(PROTO_DIR) \
	--go_out=$(STUBS_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(STUBS_DIR) --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=$(STUBS_DIR) --grpc-gateway_opt paths=source_relative \
	$(PROTO_DIR)/*/*.proto

.PHONY: clean
clean:
	$(GO) clean ./cmd/...
	rm -rf $(STUBS_DIR) 2>/dev/null
	rm -rf $(BIN_DIR)/*