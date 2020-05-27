#! /usr/bin/make -f

# Project variables.
VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECT_NAME := $(shell basename "$(PWD)")
MARKET_SERVICE := market_observer
MARKET_API := market_api
SWAGGER_API := swagger_api

# Go related variables.
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOPKG := $(.)

# Environment variables
CONFIG_FILE=$(GOBASE)/config.yml

# Go files
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# Redirect error output to a file, so we can show it in development mode.
STDERR := /tmp/.$(PROJECT_NAME)-stderr.txt

# PID file will keep the process id of the server
PID_MARKET := /tmp/.$(PROJECT_NAME).$(MARKET_SERVICE).pid
PID_MARKET_API := /tmp/.$(PROJECT_NAME).$(MARKET_API).pid
PID_SWAGGER_API := /tmp/.$(PROJECT_NAME).$(SWAGGER_API).pid
# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

all: help

## install: Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
install: go-get

## start: Start market API server, Observer, and swagger server in development mode.
start:
	@bash -c "$(MAKE) clean compile start-market-observer start-market-api start-swagger-api"

## start-market-observer: Start market observer in development mode.
start-market-observer: stop
	@echo "  >  Starting $(PROJECT_NAME) Sync"
	@-$(GOBIN)/$(MARKET_SERVICE)/market_observer -c $(CONFIG_FILE) 2>&1 & echo $$! > $(PID_MARKET)
	@cat $(PID_MARKET) | sed "/^/s/^/  \>  Sync PID: /"
	@echo "  >  Error log: $(STDERR)"

## start-market-api: Start market api in development mode.
start-market-api: stop
	@echo "  >  Starting $(PROJECT_NAME) Sync API"
	@-$(GOBIN)/$(MARKET_API)/market_api -c $(CONFIG_FILE) 2>&1 & echo $$! > $(PID_MARKET_API)
	@cat $(PID_MARKET_API) | sed "/^/s/^/  \>  Sync PID: /"
	@echo "  >  Error log: $(STDERR)"

## start-swagger-api: Start Swagger server in development mode.
start-swagger-api: stop
	@echo "  >  Starting $(PROJECT_NAME) Sync API"
	@-$(GOBIN)/$(SWAGGER_API)/swagger_api -c $(CONFIG_FILE) 2>&1 & echo $$! > $(PID_SWAGGER_API)
	@cat $(PID_SWAGGER_API) | sed "/^/s/^/  \>  Sync PID: /"
	@echo "  >  Error log: $(STDERR)"

## stop: Stop development mode.
stop:
	@-touch $(PID_MARKET) $(PID_MARKET_API) $(PID_SWAGGER_API)
	@-kill `cat $(PID_MARKET)` 2> /dev/null || true
	@-kill `cat $(PID_MARKET_API)` 2> /dev/null || true
	@-kill `cat $(PID_SWAGGER_API)` 2> /dev/null || true
	@-rm $(PID_MARKET) $(PID_MARKET_API) $(PID_SWAGGER_API)

## compile: Compile the project.
compile:
	@-touch $(STDERR)
	@-rm $(STDERR)
	@-$(MAKE) -s go-compile 2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/.*/\nError:\n/'  | sed 's/make\[.*/ /' | sed "/^/s/^/     /" 1>&2

## exec: Run given command. e.g; make exec run="go test ./..."
exec:
	GOBIN=$(GOBIN) $(run)

## clean: Clean build files. Runs `go clean` internally.
clean:
	@-rm $(GOBIN)/$(PROJECT_NAME) 2> /dev/null
	@-rm -rf mocks
	@-$(MAKE) go-clean

## generate-mocks: Creates mockfiles.
generate-mocks:
	@-$(GOBIN)/mockery -dir storage -output mocks/storage -name DB
	@-$(GOBIN)/mockery -dir storage -output mocks/storage -name ProviderList
	@-$(GOBIN)/mockery -dir market/rate -output mocks/market/rate -name RateProvider
	@-$(GOBIN)/mockery -dir market/ticker -output mocks/market/ticker -name TickerProvider
	@-$(GOBIN)/mockery -dir market/chart -output mocks/market/chart -name ChartProvider
	@-$(GOBIN)/mockery -dir services/assets -output mocks/services/assets -name AssetClient

## test: Run all unit tests.
test: go-install-mockery generate-mocks go-test

## integration: Run all integration tests.
integration: go-integration

## fmt: Run `go fmt` for all go files.
fmt: go-fmt

## govet: Run go vet.
govet: go-install-mockery generate-mocks go-vet

## golint: Run golint.
lint: go-lint-install go-lint

## docs: Generate swagger docs.
docs: go-gen-docs

## install-newman: Install Postman Newman for tests.
install-newman:
ifeq (,$(shell which newman))
	@echo "  >  Installing Postman Newman"
	@-npm install -g newman
endif

## newman: Run Postman Newman test, the host parameter is required, and you can specify the name of the test do you wanna run (transaction, token, staking, collection, domain, healthcheck, observer). e.g $ make newman test=staking host=http//localhost
newman: install-newman
	@echo "  >  Running $(test) tests"
ifeq (,$(host))
	@echo "  >  Host parameter is missing. e.g: make newman test=staking host=http://localhost:8420"
	@exit 1
endif
ifeq (,$(test))
	@bash -c "$(MAKE) newman test=healthcheck host=$(host)"
	@bash -c "$(MAKE) newman test=market host=$(host)"
else
	@newman run tests/postman/watchmarket.postman_collection.json --folder $(test) -d tests/postman/$(test)_data.json --env-var "host=$(host)"
endif

go-compile: go-get go-build

go-build:
	@echo "  >  Building market_observer binary..."
	GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/$(MARKET_SERVICE)/market_observer ./cmd/$(MARKET_SERVICE)
	@echo "  >  Building market_api binary..."
	GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/$(MARKET_API)/market_api ./cmd/$(MARKET_API)
	@echo "  >  Building swagger_api binary..."
	GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/$(SWAGGER_API)/swagger_api ./cmd/$(SWAGGER_API)

go-generate:
	@echo "  >  Generating dependency files..."
	GOBIN=$(GOBIN) go generate $(generate)

go-get:
	@echo "  >  Checking if there are any missing dependencies..."
	GOBIN=$(GOBIN) go get cmd/... $(get)

go-install:
	GOBIN=$(GOBIN) go install $(GOPKG)

go-clean:
	@echo "  >  Cleaning build cache"
	GOBIN=$(GOBIN) go clean

go-test:
	@echo "  >  Running unit tests"
	GOBIN=$(GOBIN) go test -coverprofile=coverage.txt -cover -race -v ./...

go-integration:
	@echo "  >  Running integration tests"
	GOBIN=$(GOBIN) TEST_CONFIG=$(CONFIG_FILE) go test -race -tags=integration -v ./tests/integration

go-fmt:
	@echo "  >  Format all go files"
	GOBIN=$(GOBIN) gofmt -w ${GOFMT_FILES}

go-gen-docs:
	@echo "  >  Generating swagger files"
	swag init -g ./cmd/market_api/main.go -o ./docs

go-vet:
	@echo "  >  Running go vet"
	GOBIN=$(GOBIN) go vet ./...

go-install-mockery:
	@echo "  >  Installing mockery"
	GOBIN=$(GOBIN) go get github.com/vektra/mockery/.../

go-lint-install:
	@echo "  >  Installing golint"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s

go-lint: go-install-mockery generate-mocks
	@echo "  >  Running golint"
	bin/golangci-lint

.PHONY: help

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo