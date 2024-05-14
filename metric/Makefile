GO_TEST_COVERAGE_VER ?= v2.10.1
GO_TEST_COVERAGE_CLI = $(shell command -v go-test-coverage 2> /dev/null)
GOLANGCI_LINT = $(shell command -v golangci-lint 2> /dev/null)
GOLANGCI_LINT_VERSION ?= v1.58.1

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## test: run all tests
.PHONY: test/unit
test/unit:
	go test -v -tags=unit  ./...

.PHONY: test/unit/coverage
test/unit/coverage:
	go test -count=1 -coverprofile=./coverage.out -covermode=atomic ./... -tags=unit
	@if [ -z "$(GO_TEST_COVERAGE_CLI)" ] || ! $(GO_TEST_COVERAGE_CLI) --version | grep -q "$(GO_TEST_COVERAGE_VER)$$"; then \
		echo "Installing go-test-coverage@${GO_TEST_COVERAGE_VER}"; \
		go install github.com/vladopajic/go-test-coverage/v2@${GO_TEST_COVERAGE_VER}; \
	fi
	go-test-coverage --config=./.testcoverage.yml

.PHONY: lint
lint:
	@if [ -z "$(GOLANGCI_LINT)" ] || ! $(GOLANGCI_LINT) version | grep -q "$(GOLANGCI_LINT_VERSION)$$"; then \
		echo "golangci-lint not found or version mismatch. Installing $(GOLANGCI_LINT_VERSION)..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}; \
	fi
	$(GOLANGCI_LINT) run -v
