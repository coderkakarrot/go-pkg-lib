GO_TEST_COVERAGE_VER ?= v2.10.1
GO_TEST_COVERAGE_CLI = $(shell command -v go-test-coverage 2> /dev/null)
GOLANGCI_LINT = $(shell command -v golangci-lint 2> /dev/null)
GOLANGCI_LINT_VERSION ?= v1.58.1
PROJECT_ROOT = $(shell pwd 2> /dev/null)

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## check-target-directory-option: checks target directory option provided or not
.PHONY: check-target-directory-option
check-target-directory-option:
ifndef TARGET_DIR
	$(error TARGET_DIR is not set)
endif

## tidy: tidy modfile
.PHONY: tidy
tidy: check-target-directory-option
	@echo "Running go mod tidy"
	@cd $(TARGET_DIR) && go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit: check-target-directory-option
	@echo "Running go mod verify"
	@cd $(TARGET_DIR) && go mod verify
	@echo "Running go vet ./..."
	@cd $(TARGET_DIR) && go vet ./...
	@echo "Running go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./..."
	@cd $(TARGET_DIR) && go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	@echo "Running go run golang.org/x/vuln/cmd/govulncheck@latest ./..."
	@cd $(TARGET_DIR) && go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## test/unit: run all unit tests
.PHONY: test/unit
test/unit: check-target-directory-option
	go test -v -tags=unit  ./$(TARGET_DIR)/...

## test/integration: run all integration tests
.PHONY: test/integration
test/integration: check-target-directory-option
	go test -v -tags=integration  ./$(TARGET_DIR)/...

## test/unit/coverage: checks code coverage for unit tests
.PHONY: test/unit/coverage
test/unit/coverage: check-target-directory-option
	go test -count=1 -coverprofile=./$(TARGET_DIR)/coverage.out -covermode=atomic ./$(TARGET_DIR)/... -tags=unit
	@if [ -z "$(GO_TEST_COVERAGE_CLI)" ] || ! $(GO_TEST_COVERAGE_CLI) --version | grep -q "$(GO_TEST_COVERAGE_VER)$$"; then \
		echo "Installing go-test-coverage@${GO_TEST_COVERAGE_VER}"; \
		go install github.com/vladopajic/go-test-coverage/v2@${GO_TEST_COVERAGE_VER}; \
	fi
	go-test-coverage --config=.testcoverage.yml --profile=./$(TARGET_DIR)/coverage.out

## test/integration/coverage: checks code coverage for integration tests
.PHONY: test/integration/coverage
test/integration/coverage: check-target-directory-option
	go test -count=1 -coverprofile=./$(TARGET_DIR)/coverage.out -covermode=atomic ./$(TARGET_DIR)/... -tags=integration
	@if [ -z "$(GO_TEST_COVERAGE_CLI)" ] || ! $(GO_TEST_COVERAGE_CLI) --version | grep -q "$(GO_TEST_COVERAGE_VER)$$"; then \
		echo "Installing go-test-coverage@${GO_TEST_COVERAGE_VER}"; \
		go install github.com/vladopajic/go-test-coverage/v2@${GO_TEST_COVERAGE_VER}; \
	fi
	go-test-coverage --config=.testcoverage.yml --profile=./$(TARGET_DIR)/coverage.out

## lint: run linter
.PHONY: lint
lint: check-target-directory-option
	@if [ -z "$(GOLANGCI_LINT)" ] || ! $(GOLANGCI_LINT) version | grep -q "$(GOLANGCI_LINT_VERSION)$$"; then \
		echo "golangci-lint not found or version mismatch. Installing $(GOLANGCI_LINT_VERSION)..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}; \
	fi
	@ROOT_DIR=$(pwd)
	@echo "Running golangci-lint run -v"
	@cd $(TARGET_DIR) && $(GOLANGCI_LINT) run -v --config $(PROJECT_ROOT)/.golangci.yml

## generate/unit/test/coverage: generates unit test coverage
.PHONY: generate/unit/test/coverage
generate/unit/test/coverage: check-target-directory-option
	@echo "Generating test coverage report"
	go test -count=1 -coverprofile=./$(TARGET_DIR)/coverage.out -covermode=atomic ./$(TARGET_DIR)/... -tags=unit

## generate/unit/integration/coverage: generates integration test coverage
.PHONY: generate/integration/test/coverage
generate/integration/test/coverage: check-target-directory-option
	@echo "Generating test coverage report"
	go test -count=1 -coverprofile=./$(TARGET_DIR)/coverage.out -covermode=atomic ./$(TARGET_DIR)/... -tags=integration
