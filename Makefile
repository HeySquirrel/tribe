# Inspired by https://github.com/vincentbernat/hellogopher/blob/feature/dep/Makefile

PACKAGE  = github.com/HeySquirrel/tribe
BASE     = $(GOPATH)/src/$(PACKAGE)

BIN    = $(GOPATH)/bin

TAGS = authenticated
PKGS = $(or $(PKG), $(shell cd $(BASE) && \
    env GOPATH=$(GOPATH) $(GO) list -tags=$(TAGS) ./... | grep -v "^$(PACKAGE)/vendor/"))
TESTPKGS = $(shell env GOPATH=$(GOPATH) $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))

GO      = go
GOFMT   = gofmt
TIMEOUT = 15
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

print-%  : ; @echo $* = $($*)

# Tools
GODEP = $(BIN)/dep
$(BIN)/dep: | $(BASE) ; $(info $(M) building go dep…)
	$Q go get github.com/golang/dep/cmd/dep

GOLINT = $(BIN)/golint
$(BIN)/golint: | $(BASE) ; $(info $(M) building golint…)
	$Q go get github.com/golang/lint/golint

GOCOVMERGE = $(BIN)/gocovmerge
$(BIN)/gocovmerge: | $(BASE) ; $(info $(M) building gocovmerge…)
	$Q go get github.com/wadey/gocovmerge

GOCOV = $(BIN)/gocov
$(BIN)/gocov: | $(BASE) ; $(info $(M) building gocov…)
	$Q go get github.com/axw/gocov/...

GOCOVXML = $(BIN)/gocov-xml
$(BIN)/gocov-xml: | $(BASE) ; $(info $(M) building gocov-xml…)
	$Q go get github.com/AlekSi/gocov-xml

GO2XUNIT = $(BIN)/go2xunit
$(BIN)/go2xunit: | $(BASE) ; $(info $(M) building go2xunit…)
	$Q go get github.com/tebeka/go2xunit

# Tests

TEST_DIR = $(CURDIR)/test-results
TEST_TARGETS := test-default test-bench test-short test-verbose test-race test-pr
.PHONY: $(TEST_TARGETS) test-xml check test tests
test-bench:   ARGS=-run=__absolutelynothing__ -bench=. ## Run benchmarks
test-short:   ARGS=-short        ## Run only short tests
test-verbose: ARGS=-v            ## Run tests in verbose mode with coverage reporting
test-race:    ARGS=-race         ## Run tests with race detector
test-pr:      ARGS=-tags=$(TAGS) ## Run pull request tests
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test
check test tests: fmt vendor ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests
	$Q cd $(BASE) && $(GO) test -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

test-xml: fmt vendor | $(GO2XUNIT) ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests with xUnit output
	$Q mkdir -p $(TEST_DIR)
	$Q cd $(BASE) && 2>&1 $(GO) test -timeout 20s -v $(TESTPKGS) | tee $(TEST_DIR)/tests.output
	$(GO2XUNIT) -fail -input $(TEST_DIR)/tests.output -output $(TEST_DIR)/tests.xml

COVERAGE_MODE = atomic
COVERAGE_PROFILE = $(CURDIR)/c.out
COVERAGE_XML = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML = $(COVERAGE_DIR)/index.html
.PHONY: test-coverage test-coverage-tools
test-coverage-tools: | $(GOCOVMERGE) $(GOCOV) $(GOCOVXML)
test-coverage: COVERAGE_DIR := $(TEST_DIR)/coverage
test-coverage: fmt vendor test-coverage-tools | $(BASE) ; $(info $(M) running coverage tests…) @ ## Run coverage tests
	$Q mkdir -p $(COVERAGE_DIR)/coverage
	$Q cd $(BASE) && for pkg in $(TESTPKGS); do \
		$(GO) test \
			-coverpkg=$$($(GO) list -f '{{ join .Deps "\n" }}' $$pkg | \
					grep '^$(PACKAGE)/' | grep -v '^$(PACKAGE)/vendor/' | \
					tr '\n' ',')$$pkg \
			-covermode=$(COVERAGE_MODE) \
			-coverprofile="$(COVERAGE_DIR)/coverage/`echo $$pkg | tr "/" "-"`.cover" $$pkg ;\
	 done
	$Q $(GOCOVMERGE) $(COVERAGE_DIR)/coverage/*.cover > $(COVERAGE_PROFILE)
	$Q $(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	$Q $(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)

.PHONY: lint
lint: vendor | $(GOLINT) ; $(info $(M) running golint…) @ ## Run golint
	$Q cd $(BASE) && ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GOLINT) $$pkg | tee /dev/stderr)" || ret=1 ; \
	 done ; exit $$ret

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOFMT) -l -w $$d/*.go || ret=$$? ; \
	done ; exit $$ret

# Dependency management
vendor: Gopkg.toml Gopkg.lock | $(BASE) $(GODEP) ; $(info $(M) retrieving dependencies…)
	$Q cd $(BASE) && $(GODEP) ensure
	@ln -nsf . vendor/src
	@touch $@
.PHONY: vendor-update
vendor-update: | $(BASE) $(GODEP)
ifeq "$(origin PKG)" "command line"
	$(info $(M) updating $(PKG) dependency…)
	$Q cd $(BASE) && $(GODEP) ensure -update $(PKG)
else
	$(info $(M) updating all dependencies…)
	$Q cd $(BASE) && $(GODEP) ensure -update
endif
	@ln -nsf . vendor/src
	@touch vendor


# Misc
.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf bin
	@rm -rf $(TEST_DIR)

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
