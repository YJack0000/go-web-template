include .env
export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

swag-v1: 
	swag init -g internal/controller/restful/router.go
.PHONY: swag-v1

run: swag-v1 
	 go mod tidy && go mod download && \
	 DISABLE_SWAGGER_HTTP_HANDLER='' GIN_MODE=debug CGO_ENABLED=0 go run ./cmd/main.go
.PHONY: run

build: swag-v1
	   go build -o ./tmp/main ./cmd/main.go
.PHONY: build

test:
	go test -v -cover -race ./internal/...
.PHONY: test

docker-build-dev:
	docker build -f Dockerfile.dev -t llm-dev .
.PHONY: docker-build-dev

docker-compose-dev:
	docker-compose -f docker-compose-dev.yml up
