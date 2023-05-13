GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

# Configurations are tested in Windows only, but should work relatively well in Linux with minor tweaks
ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	CONFIG_PROTO_FILES=$(shell powershell -Command "Get-ChildItem -Path .\internal\conf -Recurse -Filter "*.proto" | Resolve-Path -Relative")
	API_PROTO_FILES=$(shell powershell -Command "Get-ChildItem -Path .\api -Recurse -Filter "*.proto" | Resolve-Path -Relative")
	KEYSTORE_DIR=$(shell echo %LocalAppData%\Ethereum\Keystore)
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	CONFIG_PROTO_FILES=$(shell find internal -name ./internal/conf/*.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
	KEYSTORE_DIR=~/.ethereum/keystore
endif

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest

ETHERBASE=0x929548598a3b93362c5aa2a24de190d18e657ae0
run-geth:
	geth --datadir ../blockchain-data --keystore $(KEYSTORE_DIR) --mine --miner.etherbase $(ETHERBASE)

run:
	go run ./... -conf ./configs/config.yaml

.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal --proto_path=./third_party --go_out=paths=source_relative:./internal $(CONFIG_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api --proto_path=./third_party --go_out=paths=source_relative:./api --go-http_out=paths=source_relative:./api --go-grpc_out=paths=source_relative:./api --openapi_out=fq_schema_naming=true,default_response=false:. $(API_PROTO_FILES)

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: wire
# wire generate
wire:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...

.PHONY: all
# generate all
all:
	make api;
	make config;
	make wire;

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
