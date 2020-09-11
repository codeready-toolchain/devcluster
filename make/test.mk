COV_DIR = $(OUT_DIR)/coverage

.PHONY: test
## runs the unit tests with bundles assets
test: generate
	@echo "running the unit tests without coverage..."
	DEVCLUSTER_RESOURCE_UNIT_TEST=true DEVCLUSTER_RESOURCE_DATABASE=false go test ${V_FLAG} -race ./pkg/...

.PHONY: test-integration
## runs the integration tests
test-integration: generate
	@echo "running the integration tests without coverage..."
	DEVCLUSTER_RESOURCE_UNIT_TEST=false DEVCLUSTER_RESOURCE_DATABASE=true go test ${V_FLAG} -race ./pkg/...

.PHONY: test-all-with-coverage
## runs all the tests with coverage
test-all-with-coverage: generate
	@echo "running all the tests with coverage..."
	@-mkdir -p $(COV_DIR)
	@-rm $(COV_DIR)/coverage.txt
	DEVCLUSTER_RESOURCE_UNIT_TEST=true DEVCLUSTER_RESOURCE_DATABASE=true go test -timeout 10m -vet off ${V_FLAG} -coverprofile=$(COV_DIR)/coverage.txt -covermode=atomic ./pkg/...

.PHONY: test-all
## runs all the tests
test-all: test test-integration

.PHONY: upload-codecov-report
# Uploads the test coverage reports to codecov.io. 
# DO NOT USE LOCALLY: must only be called by OpenShift CI when processing new PR and when a PR is merged! 
upload-codecov-report: 
	# Upload coverage to codecov.io. Since we don't run on a supported CI platform (Jenkins, Travis-ci, etc.), 
	# we need to provide the PR metadata explicitely using env vars used coming from https://github.com/openshift/test-infra/blob/master/prow/jobs.md#job-environment-variables
	# 
	# Also: not using the `-F unittests` flag for now as it's temporarily disabled in the codecov UI 
	# (see https://docs.codecov.io/docs/flags#section-flags-in-the-codecov-ui)
	env
ifneq ($(PR_COMMIT), null)
	@echo "uploading test coverage report for pull-request #$(PULL_NUMBER)..."
	bash <(curl -s https://codecov.io/bash) \
		-t $(CODECOV_TOKEN) \
		-f $(COV_DIR)/coverage.txt \
		-C $(PR_COMMIT) \
		-r $(REPO_OWNER)/$(REPO_NAME) \
		-P $(PULL_NUMBER) \
		-Z
else
	@echo "uploading test coverage report after PR was merged..."
	bash <(curl -s https://codecov.io/bash) \
		-t $(CODECOV_TOKEN) \
		-f $(COV_DIR)/coverage.txt \
		-C $(BASE_COMMIT) \
		-r $(REPO_OWNER)/$(REPO_NAME) \
		-Z
endif

CODECOV_TOKEN := "8ceaf93c-f980-4cd7-8c67-7c69ae764995"
REPO_OWNER := $(shell echo $$CLONEREFS_OPTIONS | jq '.refs[0].org')
REPO_NAME := $(shell echo $$CLONEREFS_OPTIONS | jq '.refs[0].repo')
BASE_COMMIT := $(shell echo $$CLONEREFS_OPTIONS | jq '.refs[0].base_sha')
PR_COMMIT := $(shell echo $$CLONEREFS_OPTIONS | jq '.refs[0].pulls[0].sha')
PULL_NUMBER := $(shell echo $$CLONEREFS_OPTIONS | jq '.refs[0].pulls[0].number')
