## ----------------------
## Available make targets
## ----------------------
##

default: help

help: ## Display this message
	@grep -E '(^[a-zA-Z0-9_\-\.]+:.*?##.*$$)|(^##)' Makefile | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | \
	sed -e 's/\[32m##/[33m/'

##
## ----------------------
## Q.A
## ----------------------
##

qa: lint ## Run all Q.A.

LINT_CONFIG_VERSION = v1.0.2

lint:
	curl -o .golangci.yml "https://raw.githubusercontent.com/radiofrance/lint-config/${LINT_CONFIG_VERSION}/.golangci.yml"
	golangci-lint run --verbose
