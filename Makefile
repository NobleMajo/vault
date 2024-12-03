# import and use env vars if exist
-include .env.project
-include .env
export

# prepare env file and tmp dir
$(shell mkdir -p ./.tmp && touch .env)

# check required project vars
ifndef PROJECT_DISPLAY_NAME
$(error PROJECT_DISPLAY_NAME is not set)
endif

# project default vars
ifdef PROJECT_EXTRA_BUILD_ARGS
PROJECT_EXTRA_BUILD_ARGS +=
endif

PROJECT_VERSION ?= 0.0.1
PROJECT_SHORT_NAME ?= $(PROJECT_DISPLAY_NAME)
PROJECT_SHORT_NAME := $(shell echo $(PROJECT_SHORT_NAME) | sed 's/ //g' | sed 's/-/ /g' | tr '[:upper:]' '[:lower:]')
PROJECT_COMMIT_SHORT ?= $(shell echo "$$(git rev-parse --short HEAD)$$([ -n "$$(git status -s)" ] && echo -n "-modified")")
PROJECT_BUILD_ARGS ?= "$(PROJECT_EXTRA_BUILD_ARGS)-X main.Version=$(PROJECT_VERSION) -X main.Commit=$(PROJECT_COMMIT_SHORT) -X \"main.DisplayName=$(PROJECT_DISPLAY_NAME)\" -X main.ShortName=$(PROJECT_SHORT_NAME)"

# general default vars
PORT ?= 8080
HOST ?= 0.0.0.0
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOCACHE ?= $(shell if [ -d "~/.cache/go-build" ]; then realpath "~/.cache/go-build"; else mkdir -p .tmp/.cache/go-build && realpath ".tmp/.cache/go-build"; fi)


##@ These environment variables control various project configurations, including build, run, and deployment settings.
##@ They are loaded from the `.env.projects` file, which is overwritten the `.env` file variables.
##@
##@ Makefile vars
##@
##@ PROJECT_DISPLAY_NAME: projects full name, required
##@ PROJECT_VERSION: semver like project version,
##@: default is 0.0.1
##@ PROJECT_SHORT_NAME: short spaceless lowercase name,
##@: default is $PROJECT_DISPLAY_NAME
##@ PROJECT_COMMIT_SHORT: commit short hash,
##@: default is loaded via git
##@ PROJECT_BUILD_ARGS: build args,
##@: default is project base vars + $PROJECT_EXTRA_BUILD_ARGS
##@ PROJECT_EXTRA_BUILD_ARGS: additional build args,
##@: default is empty
##@
##@ Docker vars
##@
##@ DOCKER_PORT: container bind port or range (80-84),
##@: default is $PORT else 8080
##@ DOCKER_HOST: container bind host,
##@: default is $HOST else 0.0.0.0
##@ DOCKER_ARGS: container arguments for app executions,
##@: default is $ARGS
##@
##@ General vars
##@
##@ PORT: main app bind port if needed,
##@: default is empty
##@ HOST: main app bind host if needed,
##@: default is empty
##@ ARGS: arguments for app executions,
##@: default is $ARGS
##@ GOOS: target build operating system,
##@: default build system os
##@ GOARCH: target build architecture,
##@: default build system arch
##@ GOCACHE: container go cache dir,
##@: default global go cache if is dir
##@: else creates ./.tmp/.cache/go-build

##@
##@ Misc commands
##@

.PHONY: help
help: ##@ prints a command help message
	@make -s info
	@grep -F -h "##@" $(MAKEFILE_LIST) | grep -F -v grep -F | sed -e 's/\\$$//' | awk 'BEGIN {FS = ":*[[:space:]]*##@[[:space:]]*"}; \
	{ \
		if($$2 == "") \
			printf "\n"; \
		else if($$0 ~ /^#/) { \
			split($$2, arr, ": "); \
			if (length(arr) > 1) \
				printf " \033[34m%s\033[0m %s\n", arr[1], arr[2]; \
			else \
				printf "%s\n", $$2; \
		} \
		else if($$1 == "") \
			printf "%-15s^- %s\n", "", $$2; \
		else \
			printf " \033[34m%-16s\033[0m %s\n", $$1, $$2; \
	}'
	@printf "\nUsage: make <command>\n"

.PHONY: info
info: ##@ prints a project info message
	@echo "$(PROJECT_DISPLAY_NAME) version $(PROJECT_VERSION), build $(PROJECT_COMMIT_SHORT)"

.PHONY: vars
vars: ##@ prints some vars for debugging
	@echo "\nProject\n"
	@env |grep "^PROJECT_" || true
	@echo "\nDocker\n"
	@env |grep "^DOCKER_" || true
	@echo "\nGo\n"
	@env |grep "^GO" || true
	@echo "\nGeneral\n"
	@echo "ARGS: '$(ARGS)'"
	@echo "PORT: '$(PORT)'"
	@echo "HOST: '$(HOST)'"

.PHONY: env
env: ##@ prints env vars for debugging
	@env

.PHONY: clean
clean: ##@ cleans up generated files and docker cache
	@rm -fr ./.tmp ./out ./bin
	@if command -v go 2>&1 >/dev/null; then \
		echo "cleanup go..."; \
		go clean; \
		go clean -cache -fuzzcache; \
	fi
	@if command -v docker 2>&1 >/dev/null; then \
		echo "cleanup docker containers and images..."; \
		docker rm -f dev-local-$(PROJECT_SHORT_NAME)-bash > /dev/null 2>&1 \
		docker rm -f dev-local-$(PROJECT_SHORT_NAME) > /dev/null 2>&1 \
		docker image prune -f; \
	fi
	@echo "cleanup done!"
	
##@
##@ Go commands
##@

.PHONY: run
run: ##@ runs the main.go file using go run
	@go run main.go $(ARGS)

.PHONY: build
build: ##@ uses go to build the app with build args
	@touch .env
	go build \
		-ldflags=$(PROJECT_BUILD_ARGS) \
		-o bin
	@chmod +x bin

.PHONY: buildall
buildall: ##@ cross-compilation for all GOOS/GOARCH combinations
	@echo "Run test build-system build..."
	@make -s build || { echo "Test system-build build failed!"; exit 1; }
	@echo "Prepare out dir..."
	@mkdir -p out .tmp/out-bak
	@mv out/* .tmp/out-bak || true
	@echo "Start build processes..."
	@go tool dist list | while IFS=/ read -r GOOS GOARCH; do \
		( \
			go build \
				-ldflags=$(PROJECT_BUILD_ARGS) \
				-o out/$(PROJECT_SHORT_NAME)-$$GOOS-$$GOARCH && \
			chmod +x out/$(PROJECT_SHORT_NAME)-$$GOOS-$$GOARCH \
		) || { echo "Build failed for '$$GOOS'/'$$GOARCH'!"; exit 1; } \
	done && \
	wait
	@echo "All build processes finished."

.PHONY: gi
gi: ##@ installs the build binary globally
	@sudo cp ./bin /usr/local/bin/$(PROJECT_SHORT_NAME)

.PHONY: gu
gu: ##@ uninstalls the build binary globally
	@sudo rm -f /usr/local/bin/$(PROJECT_SHORT_NAME)

.PHONY: up
up: ##@ updates dependencies recursively using go get
	@echo "Update go deps recursively..."
	go get -t -u ./...

.PHONY: i
i: ##@ makes the go.mod matches the source code in the module
	@echo "Install go deps recursively..."
	go mod tidy

.PHONY: test
test: ##@ runs all GO tests recursively without coverage
	@echo "Run go tests recursively..."
	go test ./...

.PHONY: cover
cover: ##@ generates a cover.html test coverage report
	@echo "Run go tests recursively..."
	go test -coverprofile .tmp/cover.out ./...
	go tool cover -html=.tmp/cover.out -o .tmp/cover.html

.PHONY: init
init: ##@ infos, deps install, test and build
	@make -s info
	@make -s i
	@make -s test
	@make -s build

.PHONY: air
air: ##@ air watch mode and also installs air
	@go install github.com/air-verse/air@v1
	@air

##@
##@ Docker envs
##@

.PHONY: dev
dev: ##@ starts a air dev docker container
	@docker rm -f dev-local-$(PROJECT_SHORT_NAME) > /dev/null 2>&1

	docker compose run --build --rm -it --name dev-local-$(PROJECT_SHORT_NAME) -P local

.PHONY: bash
bash: ##@ runs a bash shell in the dev container
	@docker rm -f dev-local-$(PROJECT_SHORT_NAME)-bash > /dev/null 2>&1
	
	docker compose run --build --rm -it --name dev-local-$(PROJECT_SHORT_NAME)-bash --entrypoint bash -P local
