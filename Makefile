ROOT = $(shell pwd)
SERVICE_NAME = $(shell basename "$(PWD)")
SERVICE_NAME_APP = $(SERVICE_NAME)"-app"

GO ?= go
OS = $(shell uname -s | tr A-Z a-z)
export GOBIN = ${ROOT}/bin

BUILD_OUTPUT = ./bin/${SERVICE_NAME_APP}

LINT = ${GOBIN}/golangci-lint
LINT_DOWNLOAD = curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) 'latest'
GOPLANTUML = ${GOBIN}/goplantuml
GOPLANTUML_DOWNLOAD = $(GO) get github.com/jfeliu007/goplantuml/cmd/goplantuml
VERSION_TAG = $(shell git describe --tags --abbrev=0 --always)
VERSION_COMMIT = $(shell git rev-parse --short HEAD)
VERSION_DATE = $(shell git show -s --format=%cI HEAD)
VERSION = -X main.versionTag=$(VERSION_TAG) -X main.versionCommit=$(VERSION_COMMIT)  -X main.versionDate=$(VERSION_DATE) -X main.serviceName=$(SERVICE_NAME)
TPARSE = $(GOBIN)/tparse
TPARSE_DOWNLOAD = $(GO) get github.com/mfridman/tparse
COMPILEDEAMON = $(GOBIN)/CompileDaemon
COMPILEDEAMON_DOWNLOAD = $(GO) install github.com/githubnemo/CompileDaemon@latest
GOFUMPT = $(GOBIN)/gofumpt
GOFUMPT_DOWNLOAD = $(GO) install mvdan.cc/gofumpt@latest
MIGRATE_VERSION = v4.15.2
MIGRATE = ${GOBIN}/migrate
MIGRATE_DOWNLOAD = (curl --progress-bar -fL -o $(MIGRATE).tar.gz https://github.com/golang-migrate/migrate/releases/download/$(MIGRATE_VERSION)/migrate.$(OS)-amd64.tar.gz; tar -xzvf $(MIGRATE).tar.gz -C $(GOBIN); rm $(MIGRATE).tar.gz ; rm bin/LICENSE bin/README.md )
MIGRATE_CONFIG = -source file://migrations -database "${DATABASE_DRIVER}://${DATABASE_USER}:${DATABASE_PASS}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=${DATABASE_SSLMODE}"
PATH := $(PATH):$(GOBIN)

.PHONY: help
help: ## Display this help message
	@ cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build development binary file
	@ $(GO) build -ldflags '$(VERSION)' -o ${BUILD_OUTPUT} ./cmd...

.PHONY: mod
mod: ## Get dependency packages
	@ $(GO) mod tidy

.PHONY: test
test: ## Run data race detector
	@ test -e $(TPARSE) || $(TPARSE_DOWNLOAD)
	@ $(GO) test -timeout 1000s -short -race ./internal/test -json -cover | $(TPARSE) -all -smallscreen

.PHONY: test-all
test-all: ## Run data race detector
	@ test -e $(TPARSE) || $(TPARSE_DOWNLOAD)
	@ $(GO) test -timeout 1000s -short -race ./... -json -cover | $(TPARSE) -all -smallscreen

.PHONY: coverage
coverage: ## check coverage test code of sample https://penkovski.com/post/gitlab-golang-test-coverage/
	@ $(GO) test -timeout 1000s ./internal/test -coverprofile=coverage.out
	@ $(GO) tool cover -func=coverage.out
	@ $(GO) tool cover -html=coverage.out -o coverage.html;

.PHONY: uml
uml: ## Create UML diagram in diagram.puml file
	@ test -e $(GOPLANTUML) || $(GOPLANTUML_DOWNLOAD)
	@ $(GOPLANTUML) -recursive . > diagram.puml

.PHONY: migrate
migrate: ## base migrate	
	@ test -e $(MIGRATE) || $(MIGRATE_DOWNLOAD)
	@ $(MIGRATE) --version

.PHONY: migrate-create
migrate-create: ## create new migrate file	
	@ test -e $(MIGRATE) || $(MIGRATE_DOWNLOAD)	
	@ $(MIGRATE) create -ext sql -dir ./migrations -format '20060102150405' $(name)															

.PHONY: migrate-up
migrate-up:migrate ## Apply all up migrations	
	@ test -e $(MIGRATE) || $(MIGRATE_DOWNLOAD) 	
	@ $(MIGRATE) $(MIGRATE_CONFIG) up

.PHONY:	migrate-down
migrate-down:migrate ## Apply all down migrations
	@ $(MIGRATE) $(MIGRATE_CONFIG) down $(step)

.PHONY: migrate-recreate
migrate-recreate:migrate ## Apply all up migrations
	@ test -e $(MIGRATE) || $(MIGRATE_DOWNLOAD)
	@ $(MIGRATE) $(MIGRATE_CONFIG) down -all && $(MIGRATE) $(MIGRATE_CONFIG) up

.PHONY: env
env: ## create env file from .env.example and read env file & export to terminal
	@ test -e ${ROOT}/.env && echo ${ROOT}/.env exists || cp ${ROOT}/.env.example ${ROOT}/.env
	@ export $(grep -v '^#' .env | xargs -d '\n')

.PHONY: gofumpt
gofumpt: ## Lint the files
	@ test -e $(GOFUMPT) || $(GOFUMPT_DOWNLOAD)
	@ $(GOFUMPT) -l -w .
