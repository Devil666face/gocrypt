.DEFAULT_GOAL := help
PROJECT_BIN = $(shell pwd)/bin
$(shell [ -f bin ] || mkdir -p $(PROJECT_BIN))
PATH := $(PROJECT_BIN):$(PATH)
GOOS = linux
GOARCH = amd64
CGO_ENABLED = 0
VERS = $(shell git describe --tags --abbrev=0)
LDFLAGS = "-w -s -X main.vers=$(VERS)"
GCFLAGS = "-trimpath"
ASMFLAGS = "-trimpath"
# APP := $(notdir $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST))))))
ENC_BIN = encrypt
DEC_BIN = decrypt
KEY_BIN = keys
ENC_TARGET = cmd/encrypt/main.go
DEC_TARGET = cmd/decrypt/main.go
KEY_TARGET = cmd/keys/main.go

.PHONY: build \
		run \
		.install-linter \
		lint \
		.install-nil \
		nil-check \
		help

run: build test

build: ## Build release
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -asmflags=$(ASMFLAGS) -o $(PROJECT_BIN)/$(KEY_BIN) $(KEY_TARGET)
	$(KEY_BIN)
	mv id_rsa cmd/decrypt/id_rsa
	mv id_rsa.pub cmd/encrypt/id_rsa.pub
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -asmflags=$(ASMFLAGS) -o $(PROJECT_BIN)/$(ENC_BIN) $(ENC_TARGET)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -asmflags=$(ASMFLAGS) -o $(PROJECT_BIN)/$(DEC_BIN) $(DEC_TARGET)
	upx $(PROJECT_BIN)/$(ENC_BIN)
	upx $(PROJECT_BIN)/$(DEC_BIN)
	rm cmd/encrypt/id_rsa.pub cmd/decrypt/id_rsa $(PROJECT_BIN)/$(KEY_BIN)

.install-linter: ## Install linter
	[ -f golangci-lint ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) v1.54.2

lint: .install-linter ## Run linter
	golangci-lint run ./...

.install-nil: ## Install nil check
	[ -f nilaway ] || go install go.uber.org/nilaway/cmd/nilaway@latest && cp $(GOPATH)/bin/nilaway $(PROJECT_BIN)

nil-check: .install-nil ## Run nil check linter
	nilaway ./...

help:
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
