GO_VERSION=1.20
GO_DEV=./external/shell_scripts_lib/golang/go_dev.sh
GO_RELEASER=./external/shell_scripts_lib/golang/go_releaser.sh

default: help

.PHONY: update-externals
update-externals: ## Update external source dependents
	@git-deps-syncer sync-all -y

.PHONY: deps
deps: fmtcheck ## Tidy, verify and vendor go dependencies
	@${GO_DEV} deps

.PHONY: fmtcheck
fmtcheck: ## Validate Go code format and imports
	@${GO_DEV} fmt --check-only

.PHONY: fmt
fmt: ## Format Go code using gofmt style and sort imports
	@${GO_DEV} fmt

.PHONY: test
test: ## Run tests suite on host machine
	@${GO_DEV} test --dense-mode

.PHONY: test-containerized
test-containerized: ## Run tests suite within a Docker container
	@${GO_DEV} test --containerized --go-version $(GO_VERSION)

.PHONY: test-with-coverage
test-with-coverage: ## Run tests suite on host machine with coverage report
	@${GO_DEV} test --coverage

# http://localhost:9001/anchor/
.PHONY: docs-site
docs-site: ## Run a local documentation site
	@${GO_DEV} docs

# http://192.168.x.xx:9001/
.PHONY: docs-site-lan
docs-site-lan: ## Run a local documentation site (LAN available)
	@${GO_DEV} docs --lan

.PHONY: build
build: fmtcheck ## Build a binary for system OS/Arch
	@${GO_RELEASER} build 

.PHONY: build-main-package
build-main-package: fmtcheck ## Build main package for system OS/Arch
	@${GO_RELEASER} build --main-package "./cmd/anchor" 

.PHONY: delete
delete: fmtcheck ## Delete a locally installed Go binary
	@${GO_RELEASER} delete --origin gobin-local

# .PHONY: delete-remote
# delete-remote: fmtcheck ## Delete a remote GitHub release
# 	@${GO_RELEASER} delete --origin github --delete-tag ???

.PHONY: release-to-github
release-to-github: fmtcheck ## Build and publish Go binary(ies) as GitHub release
	@${GO_RELEASER} publish \
		--main-package "./cmd/anchor" \
		--release-type github \
		--release-tag $(shell cat ./resources/version.txt)

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


