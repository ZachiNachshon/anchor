TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: help

.PHONY: setup
setup: fmtcheck ## Update dependencies, preparing this repo pre-build
	go mod tidy
	go mod verify
	go mod vendor

.PHONY: build
build: fmtcheck ## Build binary with format check (destination: PWD)
	go build -o ./anchor ./cmd/anchor/*.go

.PHONY: build-to-gopath
build-to-gopath: fmtcheck ## Build binary with format check (destination: GOPATH/bin)
	go build -o $(GOPATH)/bin/anchor ./cmd/anchor/*.go

.PHONY: build-ci
build-ci: fmtcheck ## Build binary with format check (destination: GOBIN on CI)
	go build ./...
	go install ./cmd/anchor/*.go

.PHONY: fmt
fmt: ## Format go files
	gofmt -w $(GOFMT_FILES)

.PHONY: fmtcheck
fmtcheck: ## Check go files format validity
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: run-tests
run-tests: ## Run tests suite locally
	@go test -v $(TEST) -json -cover -covermode=count -coverprofile=coverage.out.temp | tparse -all -notests
	@cat coverage.out.temp | grep -v '_testkit\|_fakes' > coverage.out
	@go tool cover -func coverage.out | grep total | awk '{print $3}'
	@# go test -v $(TEST) -json -cover | tparse -all -notests
	@# go test -v $(TEST) -cover time config/*.go

.PHONY: run-tests-ci
run-tests-ci: ## Run tests suite on CI containerized environment
	@go test -v $(TEST) -json -cover -covermode=count -coverprofile=coverage.out.temp | tparse -all -top
	@cat coverage.out.temp | grep -v '_testkit\|_fakes' > coverage.out
	@# -coverprofile=coverage.out was added for GitHub workflow integration with jandelgado/gcov2lcov-action
	@# Error:
	@#   /tmp/gcov2lcov-linux-amd64 -infile coverage.out -outfile coverage.lcov
	@#   2021/08/01 07:21:57 error opening input file: open coverage.out: no such file or directory

.PHONY: run-tests-containerized
run-tests-containerized: ## Run tests suite locally containerized
	docker run -it \
        -v $(PWD):/home/anchor \
        --entrypoint /bin/sh golang:1.16 \
        -c 'go get github.com/mfridman/tparse; cd /home/anchor; make run-tests-ci'
#        -c 'cd /home/anchor; go test -v $(TEST) -json -cover'

.PHONY: release
release: fmtcheck ## Create release artifacts in GitHub with version from resources/version.txt
	@sh -c "'$(CURDIR)/scripts/release.sh'"

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


