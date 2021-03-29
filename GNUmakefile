GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=rancher2
TEST?="./${PKG_NAME}"
PROVIDER_NAME=terraform-provider-rancher2

default: build

build: validate
	@sh -c "'$(CURDIR)/scripts/gobuild.sh'"

docker-build: 
	@sh -c "'$(CURDIR)/scripts/gobuild_docker.sh'"

build-rancher: validate-rancher
	@sh -c "'$(CURDIR)/scripts/gobuild.sh'"

validate-rancher: validate test

validate: fmtcheck lint vet

package-rancher: 
	@sh -c "'$(CURDIR)/scripts/gopackage.sh'"

test: fmtcheck
	@echo "==> Running testing..."
	go test $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: 
	@sh -c "'$(CURDIR)/scripts/gotestacc.sh'"

docker-testacc: 
	@sh -c "'$(CURDIR)/scripts/gotestacc_docker.sh'"

upgrade-rancher: 
	@sh -c "'$(CURDIR)/scripts/start_rancher.sh'"

vet:
	@echo "==> Checking that code complies with go vet requirements..."
	@go vet $$(go list ./... | grep -v vendor/); if [ $$? -gt 0 ]; then \
		echo ""; \
		echo "WARNING!! Expected vet reported construct:"; \
		echo "rancher2/schema_secret_v2.go:20:2: struct field Type repeats json tag \"type\" also at ../../../../github.com/rancher/norman@v0.0.0-20210225010917-c7fd1e24145b/types/types.go:66"; \
		echo "";\
		echo "If vet reported more suspicious constructs, please check and"; \
		echo "fix them if necessary, before submitting the code for review."; \
	fi

lint:
	@echo "==> Checking that code complies with golint requirements..."
	@GO111MODULE=off go get -u golang.org/x/lint/golint
	@if [ -n "$$(golint $$(go list ./...) | grep -v 'should have comment.*or be unexported' | tee /dev/stderr)" ]; then \
		echo ""; \
		echo "golint found style issues. Please check the reported issues"; \
		echo "and fix them if necessary before submitting the code for review."; \
    	exit 1; \
	fi

bin:
	go build -o $(PROVIDER_NAME)

fmt:
	gofmt -s -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor:
	@echo "==> Updating vendor modules..."
	@GO111MODULE=on go mod vendor

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile bin vendor
