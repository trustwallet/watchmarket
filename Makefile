#! /usr/bin/make -f

# Project variables.
VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECT_NAME := $(shell basename "$(PWD)")
MARKET_SERVICE := worker
MARKET_API := api
MARKET_SEED_DB := seed
MARKET_PROXY := proxy
MARKET_PG_HEALTH := pg-health

# Go related variables.
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOPKG := $(.)
# A valid GOPATH is required to use the `go get` command.
# If $GOPATH is not specified, $HOME/go will be used by default
GOPATH := $(if $(GOPATH),$(GOPATH),~/go)

DOCKER_REDIS_IMAGE_NAME := redis
DOCKER_LOCAL_DB_IMAGE_NAME := test_db
DOCKER_LOCAL_DB_USER :=user
DOCKER_LOCAL_DB_PASS :=pass
DOCKER_LOCAL_DB := watchmarket

DOCKER_REPOSITORY := trust/watchmarket
HASH ?= local

# Environment variables
CONFIG_FILE=config.yml

# Go files
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# Redirect error output to a file, so we can show it in development mode.
STDERR := /tmp/.$(PROJECT_NAME)-stderr.txt

# PID file will keep the process id of the server
PID_MARKET := /tmp/.$(PROJECT_NAME).$(MARKET_SERVICE).pid
PID_MARKET_API := /tmp/.$(PROJECT_NAME).$(MARKET_API).pid

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

all: help

## install: Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
install: go-get

## start: Start market API server, Observer, and swagger server in development mode.
start:
	@bash -c "$(MAKE) clean compile start-market-observer start-market-api"

## start-market-observer: Start market observer in development mode.
start-market-observer: stop
	@echo "  >  Starting $(PROJECT_NAME) Sync"
	@-$(GOBIN)/$(MARKET_SERVICE)/worker -c $(CONFIG_FILE) 2>&1 & echo $$! > $(PID_MARKET)
	@cat $(PID_MARKET) | sed "/^/s/^/  \>  Sync PID: /"
	@echo "  >  Error log: $(STDERR)"

## start-market-api: Start market api in development mode.
start-market-api: stop
	@echo "  >  Starting $(PROJECT_NAME) Sync API"
	@-$(GOBIN)/$(MARKET_API)/api -c $(CONFIG_FILE) 2>&1 & echo $$! > $(PID_MARKET_API)
	@cat $(PID_MARKET_API) | sed "/^/s/^/  \>  Sync PID: /"
	@echo "  >  Error log: $(STDERR)"

## stop: Stop development mode.
stop:
	@-touch $(PID_MARKET) $(PID_MARKET_API)
	@-kill `cat $(PID_MARKET)` 2> /dev/null || true
	@-kill `cat $(PID_MARKET_API)` 2> /dev/null || true
	@-rm $(PID_MARKET) $(PID_MARKET_API)

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

## test: Run all unit tests.
test: go-test

## integration: Run all integration tests.
integration: go-integration

## fmt: Run `go fmt` for all go files.
fmt: go-fmt

## govet: Run go vet.
govet: generate-mocks go-vet

## golint: Run golint.
lint: go-lint-install go-lint

## docs: Generate swagger docs.
docs: go-gen-docs

docker-shutdown:
	@echo "  >  Shutdown docker containers..."
	@-bash -c "docker rm -f $(DOCKER_LOCAL_DB_IMAGE_NAME) 2> /dev/null"
	@-bash -c "docker rm -f $(DOCKER_REDIS_IMAGE_NAME) 2> /dev/null"

start-docker-services: docker-shutdown
	docker run -d -p 5432:5432 --name $(DOCKER_LOCAL_DB_IMAGE_NAME) -e POSTGRES_USER=$(DOCKER_LOCAL_DB_USER) -e POSTGRES_PASSWORD=$(DOCKER_LOCAL_DB_PASS) -e POSTGRES_DB=$(DOCKER_LOCAL_DB) postgres
	docker run -d -p 6379:6379 --name $(DOCKER_REDIS_IMAGE_NAME) redis

seed-db:
	@echo "  >  Seeding db"
	sleep 1
	docker cp seed/. $(DOCKER_LOCAL_DB_IMAGE_NAME):/docker-entrypoint-initdb.d/

	@echo "  >  Seeding watchmarket_public_tickers"
	docker exec -it $(DOCKER_LOCAL_DB_IMAGE_NAME) psql -U $(DOCKER_LOCAL_DB_USER) -d $(DOCKER_LOCAL_DB) -f /docker-entrypoint-initdb.d/watchmarket_public_tickers.sql > /dev/null 2>&1
	@echo "  >  Seeding watchmarket_public_rates"
	docker exec -it $(DOCKER_LOCAL_DB_IMAGE_NAME) psql -U $(DOCKER_LOCAL_DB_USER) -d $(DOCKER_LOCAL_DB) -f /docker-entrypoint-initdb.d/watchmarket_public_rates.sql > /dev/null 2>&1

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
	@echo "  >  Building worker binary..."
	GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/$(MARKET_SERVICE)/worker ./cmd/$(MARKET_SERVICE)
	@echo "  >  Building api binary..."
	GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/$(MARKET_API)/api ./cmd/$(MARKET_API)

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
	GOBIN=$(GOBIN) TEST_CONFIG=$(CONFIG_FILE) go test -race -tags=integration -v ./tests/integration/...

go-fmt:
	@echo "  >  Format all go files"
	GOBIN=$(GOBIN) gofmt -w ${GOFMT_FILES}

install-swag:
ifeq (,$(wildcard test -f $(GOPATH)/bin/swag))
	@echo "  >  Installing swagger"
	@-bash -c "go get github.com/swaggo/swag/cmd/swag"
endif

swag: install-swag
	@bash -c "$(GOPATH)/bin/swag init --parseDependency -g ./cmd/api/main.go -o ./docs"

go-vet:
	@echo "  >  Running go vet"
	GOBIN=$(GOBIN) go vet ./...

go-lint-install:
ifeq (,$(wildcard test -f bin/golangci-lint))
	@echo "  >  Installing golint"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s
endif

go-lint:
	@echo "  >  Running golint"
	bin/golangci-lint run --timeout=2m

.PHONY: help

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo