BUILD_TARGETS += openapi api-docs
TEST_TARGETS += api-lint
CLEAN_TARGETS += api-clean

REDOC_VERSION ?= v2.0.0-rc.18
REDOC_CLI_VERSION ?= 0.9.6-1
PRISM_VERSION ?= 3.1.1
SPECTRAL_VERSION ?= 5.3.0
SPECCY_VERSION ?= 0.11.0
YAMLLINT_VERSION ?= 1.20
 # no tags in dockerhub for this one, but latest in github is v0.2.8
OPENAPI_SPEC_VALIDATOR_VERSION ?= latest
OPENAPI_GENERATOR_VERSION ?= v4.3.0
MINISPEC_VERSION ?= v0.2.0

_docker_opts := --user $(shell id -u):$(shell id -g) --rm
_docker_opts += --volume $(CURDIR):/local

MINISPEC = docker run $(_docker_opts) confluent-docker.jfrog.io/confluentinc/minispec:$(MINISPEC_VERSION)
MINISPEC_ARGS ?= -vvv

YAMLLINT = docker run $(_docker_opts) cytopia/yamllint:$(YAMLLINT_VERSION)
YAMLLINT_CONF ?= .yamllint

OPENAPI_SPEC_VALIDATOR = docker run $(_docker_opts) p1c2u/openapi-spec-validator:$(OPENAPI_SPEC_VALIDATOR_VERSION)

SPECCY = docker run $(_docker_opts) wework/speccy:$(SPECCY_VERSION)

SPECTRAL = docker run $(_docker_opts) stoplight/spectral:$(SPECTRAL_VERSION)

REDOC = docker run $(_docker_opts) confluent-docker.jfrog.io/confluentinc/redoc-cli:$(REDOC_CLI_VERSION)
REDOC_OPTIONS ?= -t /local/mk-include/resources/redoc.hbs --templateOptions.segmentWriteKey "$(REDOC_SEGMENT_KEY)"
REDOC_TARGETS ?= mk-include/resources/redoc.hbs
API_DOC_TITLE ?= API Reference Documentation

PRISM = docker run $(_docker_opts) stoplight/prism:$(PRISM_VERSION)

OPENAPI_GENERATOR = docker run $(_docker_opts) openapitools/openapi-generator-cli:$(OPENAPI_GENERATOR_VERSION)

.SECONDEXPANSION:

 # $(H) helper to hide the command being run using @, except on CI. Usage: $(H)$(CMD)...
ifneq ($(CI),true)
H := @
endif

 # $(OPEN) helper to automatically open a window
ifeq ($(HOST_OS),linux)
OPEN := xdg-open
else
OPEN := open
endif
ifeq ($(shell which $(OPEN)),)
OPEN := @echo Browse to
endif

 # Note: divider blocks and non-doc comments must be indented
 # to avoid showing up in `mmake help` output: https://github.com/tj/mmake

.PHONY: api-resources
api-resources:
	cp mk-include/resources/.yamllint ./.yamllint

.PHONY: api-clean
## Clean all generated artifacts that aren't maintained in Git
api-clean:
	rm -rf $(addsuffix /sdk,$(API_SPEC_DIRS))
	rm -rf $(addsuffix /postman.json,$(API_SPEC_DIRS))
	rm -rf $(addsuffix /loadtest,$(API_SPEC_DIRS))

 #############
 ## OPENAPI ##
 #############

.PHONY: openapi
## Generate OpenAPI from Minispec for all APIs
openapi: $$(addsuffix /openapi.yaml,$$(API_SPEC_DIRS))

.PHONY: api-spec
## Generate OpenAPI from Minispec for all APIs (alias: 'openapi')
api-spec: openapi

## Generate OpenAPI from Minispec for an API
%/openapi.yaml: %/minispec.yaml
	$(H)$(MINISPEC) /local/$< $(MINISPEC_ARGS) --out /local/$@

.PHONY: api-lint
## Lint all API specifications with all linters
api-lint: api-lint-yaml api-lint-openapi

.PHONY: api-lint-openapi
## Lint the OpenAPI spec with all linters
api-lint-openapi: api-lint-openapi-spec-validator api-lint-spectral

.PHONY: api-lint-openapi-spec-validator
## Lint the OpenAPI spec using openapi-spec-validator
api-lint-openapi-spec-validator: $$(addsuffix /openapi-spec-validator,$$(API_SPEC_DIRS))

.PHONY: api-lint-yaml
## Lint the OpenAPI spec using yamllint
api-lint-yaml: $$(addsuffix /yamllint,$$(API_SPEC_DIRS))

.PHONY: api-lint-spectral
## Lint the OpenAPI spec using spectral
api-lint-spectral: $$(addsuffix /spectral,$$(API_SPEC_DIRS))

.PHONY: api-lint-speccy
## (POC) Lint the OpenAPI spec using speccy
api-lint-speccy: $$(addsuffix /speccy,$$(API_SPEC_DIRS))

.PHONY: %/yamllint
## Lint the API against yamllint
%/yamllint:
ifneq ($(wildcard $(YAMLLINT_CONF)),)
	$(YAMLLINT) \
		-f colored \
		-c /local/$(YAMLLINT_CONF) \
		/local/$*
else
	$(warning Create a $(YAMLLINT_CONF) file to enable YAML linting of OpenAPI spec)
endif

.PHONY: %/openapi-spec-validator
## Lint the API against openapi-spec-validation
%/openapi-spec-validator:
	$(OPENAPI_SPEC_VALIDATOR) "/local/$*/openapi.yaml"

.PHONY: %/spectral
## Lint the API against Spectral rules
%/spectral: %/openapi.yaml
	$(SPECTRAL) lint "/local/$*/openapi.yaml"

.PHONY: %/speccy
## Lint the API against Speccy rules
## Create a speccy.yaml file in the project root to customize config
%/speccy: %/openapi.yaml
	$(SPECCY) lint /local/$<

 ###########
 ## REDOC ##
 ###########

.PHONY: api-docs
## Generate ReDoc docs for all APIs
api-docs: $$(addsuffix /openapi.html,$$(API_SPEC_DIRS))

## Generate API docs from OpenAPI using ReDoc
## Example: make ccloud/openapi.html API_DOC_TITLE="My HTML Page Title" REDOC_SEGMENT_KEY=ABC123
%/openapi.html: %/openapi.yaml $(REDOC_TARGETS)
	$(REDOC) bundle /local/$< \
		--output /local/$@ \
		--title "$(API_DOC_TITLE)" \
		--options.theme.colors.primary.main='#0074A2' \
		--options.theme.colors.primary.light='#00AFBA' \
		--options.theme.colors.primary.dark='#173361' \
		--cdn \
		$(REDOC_OPTIONS)
ifeq ($(CI),true)
ifeq ($(BRANCH_NAME),$(MASTER_BRANCH))
	artifact push project --force $@
else
	artifact push workflow $@
endif
endif

.PHONY: api-redoc-serve
## Serve the ReDoc docs
## This is useful for viewing the docs while iterating on spec development
api-redoc-serve: $$(addsuffix /redoc-serve,$$(lastword $$(API_SPEC_DIRS)))

.PHONY: %/redoc-serve
## Serve the ReDoc docs for an API (in the foreground)
## This is useful for viewing the docs while iterating on spec development
%/redoc-serve: _docker_opts += --init --publish 8080:8080
%/redoc-serve: %/openapi.yaml
	@sleep 3 && $(OPEN) http://localhost:8080
	$(REDOC) serve /local/$< --watch

.PHONY: %/redoc-start
## Start the ReDoc server for an API (in the background)
## This is useful for viewing the docs while iterating on spec development
%/redoc-start: _docker_opts += --detach --publish 8080:8080
%/redoc-start: %/openapi.yaml
	@sleep 3 && $(OPEN) http://localhost:8080
	$(REDOC) serve /local/$<

.PHONY: %/redoc-stop
## Stop the ReDoc server for an API (in the background)
%/redoc-stop:
	docker stop $(shell docker ps -q -f "ancestor=confluent-docker.jfrog.io/confluentinc/redoc-cli")

 ################
 ## MOCK (POC) ##
 ################

.PHONY: api-mock
## (POC) Run the Prism mock server for an API
api-mock: $$(addsuffix /prism,$$(lastword $$(API_SPEC_DIRS)))

.PHONY: %/prism
## (POC) Run the Prism mock server for an API
%/prism: _docker_opts += --init --publish 4010:4010
%/prism:
	$(PRISM) mock -h 0.0.0.0 "/local/$*/openapi.yaml"

 ################
 ## SDKs (POC) ##
 ################

.PHONY: sdk
## (POC) Generate SDKs for all APIs in all languages
sdk: sdk-go sdk-java

.PHONY: sdk-go
## (POC) Generate SDKs for all APIs in Golang
sdk-go: $$(addsuffix /sdk/go,$$(API_SPEC_DIRS))

.PHONY: sdk-java
## (POC) Generate SDKs for all APIs in Java
sdk-java: $$(addsuffix /sdk/java,$$(API_SPEC_DIRS))

## (POC) Generate Golang SDKs for an API
%/sdk/go: %/openapi.yaml
	$(OPENAPI_GENERATOR) generate -g go \
		-i /local/$< -o /local/$@ \
		--package-name v1

## (POC) Generate Java SDKs for an API
%/sdk/java: %/openapi.yaml
	$(OPENAPI_GENERATOR) generate -g java \
		-i /local/$< -o /local/$@ \
		--api-package io.confluent.$*.api \
		--model-package io.confluent.$*.model \
		--invoker-package io.confluent.$*.client \
		--group-id io.confluent --artifact-id $*-java-client

 ###############################
 ## Postman collections (POC) ##
 ###############################

.PHONY: api-postman
## (POC) Generate Postman collection for all APIs
api-postman: $$(addsuffix /postman.json,$$(API_SPEC_DIRS))

## (POC) Generate Postman collection for an API
%/postman.json: %/openapi.yaml
	cd postman && npm install && \
		node node_modules/openapi-to-postmanv2/bin/openapi2postmanv2.js -s ../$*/openapi.yaml -o ../$*/postman.json

 ################################
 ## Load Test in Gatling (POC) ##
 ################################

.PHONY: api-loadtest
## (POC) Generate Gatling load test for all APIs
api-loadtest: $$(addsuffix /loadtest,$$(API_SPEC_DIRS))

## (POC) Generate Gatling load test for an API
%/loadtest: %/openapi.yaml
	$(OPENAPI_GENERATOR) generate -i /local/$*/openapi.yaml -g scala-gatling -o /local/$*/loadtest
