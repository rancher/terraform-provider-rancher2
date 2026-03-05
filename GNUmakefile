default: fmt lint build install generate test testlong

fmt:
	gofmt -s -w -e .

lint:
	golangci-lint run

build:
	rm -f ./bin/terraform-provider-rancher2
	go build -o ./bin/ -v ./...

install:
	go install -v ./...

generate:
	cd tools; go generate ./...

test: # run unit tests
	gotestsum --format standard-verbose --jsonfile report.json --post-run-command "./test/unit/summarize.sh" -- ./internal/provider/... -v -p=10 -timeout=300s -cover

startlocal:
	./test/short/scripts/local_rancher.sh

stoplocal:
	./test/short/scripts/local_rancher.sh cleanup

ts: # run short tests without launching Rancher
	cd ./test/short; go test -count=1 -run='TestAccDevResource' -v; cd ../../;

testshort: startlocal ts stoplocal # run short acceptance tests

testlong: # run e2e tests
	./test/long/scripts/run_tests.sh

dt: # run specific unit test eg. `make dt -- t=<testname>`
	gotestsum --format standard-verbose -- ./... -v -run "$(t)"

et: build # run specific acceptance test eg. `make et -- t=<testname>`
	./test/long/scripts/run_tests.sh -t $(t)

clean: # clean up test leftovers eg. `make clean -- i=<identifier>`
	./test/long/scripts/run_tests.sh -c $(i)

.PHONY: fmt lint build install generate test testlong testshort dt et clean startlocal stoplocal

