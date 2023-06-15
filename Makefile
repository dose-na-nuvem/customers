.DEFAULT_GOAL := all

.PHONY: all
all: lint test build

.PHONY: lint
lint:
	@golangci-lint --timeout 120s run ./...

.PHONY: test
test:
	@go test -v ./...

.PHONY: build
build:
	@go build -o ./_build/customers .

.PHONY: air
air:
	@air

.PHONY: install-tools
install-tools:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.2
	@curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v23.2/protoc-23.2-linux-x86_64.zip
	@unzip protoc-23.2-linux-x86_64.zip -d ${HOME}/.local
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

.PHONY: protoc
protoc:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/customer/customer.proto 
