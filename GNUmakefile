TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: build

build: fmtcheck
	go build
	go install

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

release: fmtcheck
	@sh -c "'$(CURDIR)/scripts/gorelease.sh'"

.PHONY: build fmt fmtcheck release


