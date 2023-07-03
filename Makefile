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
	@if ! air version &> /dev/null; then go install github.com/cosmtrek/air@v1.44.0; fi
	@air

.PHONY: install-tools
install-tools:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.2
	@curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v23.2/protoc-23.2-linux-x86_64.zip
	@unzip protoc-23.2-linux-x86_64.zip -d ${HOME}/.local
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3
	@go install github.com/cloudflare/cfssl/cmd/cfssl@v1.6.4
	@go install github.com/cloudflare/cfssl/cmd/cfssljson@v1.6.4

.PHONY: cert
cert:
	@mkdir -p ./tmp
	@mkdir -p ./local/certs
	@echo '{ "hosts": [ "localhost", "127.0.0.1" ], "key": { "algo": "rsa", "size": 2048 }, "names": [ { "O": "Customers (Projeto Dose na Nuvem)" } ]}' > ./tmp/ca-csr.json
	@cfssl genkey -initca ./tmp/ca-csr.json | cfssljson -bare ./local/certs/ca
	@echo '{ "hosts": [ "localhost", "127.0.0.1" ], "key": { "algo": "rsa", "size": 2048 }, "names": [ { "O": "Customers (Projeto Dose na Nuvem)" } ]}' > ./tmp/cert-csr.json
	@cfssl gencert -ca ./local/certs/ca.pem -ca-key ./local/certs/ca-key.pem ./tmp/cert-csr.json | cfssljson -bare ./local/certs/cert

.PHONY: protoc
protoc:
	@${HOME}/.local/bin/protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/customer/customer.proto
