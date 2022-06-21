BINARY_NAME=anchor
GO_VERSION=1.8
TESTS_PATH?=./...

default: help

.PHONY: align-deps
align-deps: fmtcheck ## Update dependencies, preparing this repo pre-build
	@sh -c "'$(CURDIR)/external/shell-scripts-lib/golang/deps.sh'"

.PHONY: build
build: fmtcheck ## Build binary for current OS/Arch (destination: PWD)
	@sh -c "'$(CURDIR)/external/shell-scripts-lib/golang/build.sh' \
		'action: build' \
		'binary_name: $(BINARY_NAME)' \
		'dist_path: dist' \
		'go_files_path: ./cmd/$(BINARY_NAME)'"

.PHONY: build-to-gopath
build-to-gopath: fmtcheck ## Build binary with format check (destination: GOPATH/bin)
	@sh -c "'$(CURDIR)/external/shell-scripts-lib/golang/build.sh' \
		'action: build' \
		'binary_name: $(BINARY_NAME)' \
		'dist_path: $(GOPATH)/bin' \
		'go_files_path: ./cmd/$(BINARY_NAME)'"

.PHONY: build-os-arch
build-os-arch: fmtcheck ## Build binaries for specific OS/Arch (destination: PWD/dist)
	@sh -c "'$(CURDIR)/external/shell-scripts-lib/golang/build.sh' \
		'action: build-os-arch' \
		'binary_name: $(BINARY_NAME)' \
		'dist_path: dist' \
		'go_files_path: ./cmd/$(BINARY_NAME)'"

.PHONY: build-to-gobin-ci
build-to-gobin-ci: fmtcheck ## Build binary with format check (destination: GOBIN on CI)
	go build ./...
	go install ./cmd/anchor/*.go

.PHONY: fmt
fmt: ## Format go files
	@sh -c "'$(CURDIR)/external/shell-scripts-lib/golang/fmt.sh' --format"

.PHONY: fmtcheck
fmtcheck: ## Check go files format validity
	@sh -c "'$(CURDIR)/external/shell-scripts-lib/golang/fmt.sh' --check"

.PHONY: test
test: ## Run tests suite locally
	@sh -c "'$(CURDIR)/external/shell-scripts-lib/golang/test.sh' \
		'action: local' \
		'tests_path: $(TESTS_PATH)'"

.PHONY: test-github-ci
test-github-ci: ## Run tests suite on CI containerized environment
	@sh -c "'$(CURDIR)/external/shell-scripts-lib/golang/test.sh' \
		'action: github-ci' \
		'tests_path: $(TESTS_PATH)'"

.PHONY: test-containerized
test-containerized: ## Run tests suite locally containerized
	@sh -c "'$(CURDIR)/external/shell-scripts-lib/golang/test.sh' \
		'action: containerized' \
		'tests_path: $(TESTS_PATH)' \
		'project_root_path: $(PWD)' \
		'binary_name: $(BINARY_NAME)' \
		'go_version: $(GO_VERSION)'"

.PHONY: release-version
release-version: fmtcheck ## Create release artifacts in GitHub with version from resources/version.txt
	@sh -c "'$(CURDIR)/external/shell-scripts-lib/golang/goreleaser.sh' \
		'config_file_path: ./.goreleaser.yml' \
		'version_file_path: ./resources/version.txt'"

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


