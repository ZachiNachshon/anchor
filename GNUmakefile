TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: help

build: fmtcheck ## Build project, format check and install to bin folder
	go build
	go install

.PHONY: run-tests-ci
run-tests-ci: ## Run tests suites on CI containerized environment
	@go test -v $(TEST) -json -cover -covermode=count -coverprofile=coverage.out 	| tparse -all -top
	# -coverprofile=coverage.out was added for GitHub workflow integration with jandelgado/gcov2lcov-action
	# Error:
	#   /tmp/gcov2lcov-linux-amd64 -infile coverage.out -outfile coverage.lcov
    #   2021/08/01 07:21:57 error opening input file: open coverage.out: no such file or directory
	# Ignore:
	#   go test -v $(TEST) -json -cover | tparse -all -notests
	#   go test -v $(TEST) -cover time config/*.go

.PHONY: run-tests-dockerized
run-tests-dockerized: ## Run tests suites dockerized
	docker run -it \
        -v $(PWD):/home/anchor \
        --entrypoint /bin/sh golang:1.14.15 \
        -c 'cd /home/anchor; make run-tests-ci'
#        -c 'cd /home/anchor; go test -v $(TEST) -json -cover'

fmt: ## Format go files
	gofmt -w $(GOFMT_FILES)

fmtcheck: ## Check go files format validity
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

release: fmtcheck ## Create release artifacts with format check
	@sh -c "'$(CURDIR)/scripts/gorelease.sh'"

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


