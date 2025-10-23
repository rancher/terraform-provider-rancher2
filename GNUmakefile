default: fmt lint build install generate test testacc

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
	gotestsum --format standard-verbose --jsonfile report.json --post-run-command "./test/summarize.sh" -- ./... -v -p=10 -timeout=300s -cover

testacc: # run all acceptance tests
	./run_tests.sh

dt: # run specific unit test eg. `make dt -- t=<testname>`
	gotestsum --format standard-verbose -- $(t)

et: build # run specific acceptance test eg. `make et -- t=<testname>`
	./run_tests.sh -t $(t)

clean: # clean up test leftovers eg. `make clean -- i=<identifier>`
	./run_tests.sh -c $(i)

.PHONY: fmt lint build install generate test testacc dt et clean
