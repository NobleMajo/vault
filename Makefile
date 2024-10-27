DISPLAY_NAME := WireGaurdGenerator
SHORT_NAME := wgg
VERSION := 1.0.0

COMMIT := $(shell git rev-parse --short HEAD)
BUILD_ARGS := "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.DisplayName=$(DISPLAY_NAME) -X main.ShortName=$(SHORT_NAME)"
PORT ?= 8080

-include .env
export

## info: prints a project info message
.PHONY: info
info:
	@echo "$(DISPLAY_NAME) version $(VERSION), build $(COMMIT)"

## run: uses go to start the main.go
.PHONY: run
run:
	@go run main.go

## build: uses go to build the app with build args
.PHONY: build
build:
	@touch .env
	go build \
		-ldflags=$(BUILD_ARGS) \
		-o bin
	chmod +x bin

## clean: cleans up the tmp, build and docker cache
.PHONY: clean
clean:
	@rm -f bin
	@rm -fr ./tmp
	@if command -v go 2>&1 >/dev/null; then \
		echo "cleanup go..."; \
		go clean; \
		go clean -cache -fuzzcache; \
	fi
	@if command -v docker 2>&1 >/dev/null; then \
		echo "cleanup docker..."; \
		CACHE_DIR="" PORT="" docker compose down --remove-orphans --rmi all; \
		docker image prune -f; \
	fi
	@echo "cleanup done!"
	@echo "WARNING: the .env file still exists!"


## update: updates dependencies
.PHONY: update
update:
	go get -t -u ./...

## test: runs all tests without coverage
.PHONY: test
test:
	go test ./...

## init: prepares ands builds
.PHONY: init
init:
	@touch .env
	@echo "update deps..."
	@go mod tidy
	@echo "testing..."
	@make -s test
	@echo "building..."
	@make -s build

## air: starts the go bin in air watch mode
.PHONY: air
air:
	@go install github.com/air-verse/air@v1
	@air

## dev: starts a dev docker container
.PHONY: dev
dev:
	@touch .env
	$(eval CACHE_DIR = .tmp/.cache/go-build)
	@if [ -d ~/.cache/go-build ]; then \
		$(eval CACHE_DIR = ~/.cache/go-build) \
		echo "Use users go-build cache dir."; \
	else \
		mkdir -p $(CACHE_DIR); \
		echo "Use local go-build cache dir."; \
	fi
	
	@docker rm -f $(SHORT_NAME)-local-dev > /dev/null 2>&1

	PORT=${PORT} \
		CACHE_DIR=${CACHE_DIR} \
		docker compose run --build --rm -it --name $(SHORT_NAME)-local-dev -P local

## exec: starts a bash in a dev container
.PHONY: exec
exec:
	@touch .env
	$(eval CACHE_DIR = .tmp/.cache/go-build)
	@if [ -d ~/.cache/go-build ]; then \
		$(eval CACHE_DIR = ~/.cache/go-build) \
		echo "Use users go-build cache dir."; \
	else \
		mkdir -p $(CACHE_DIR); \
		echo "Use local go-build cache dir."; \
	fi

	@docker rm -f $(SHORT_NAME)-local-bash > /dev/null 2>&1
	
	PORT=${PORT} \
		CACHE_DIR=${CACHE_DIR} \
		docker compose run --build --rm -it --name $(SHORT_NAME)-local-bash --entrypoint bash -P local