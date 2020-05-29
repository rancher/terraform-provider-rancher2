# Set shell to bash
SHELL := /bin/bash

# Use this variable to specify a different make utility (e.g. remake --profile)
MAKE ?= make

# Include this file first
_empty :=
_space := $(_empty) $(empty)

# Master branch
MASTER_BRANCH ?= master

RELEASE_TARGETS += $(_empty)
BUILD_TARGETS += $(_empty)
TEST_TARGETS += $(_empty)
CLEAN_TARGETS += $(_empty)

# If this variable is set, release will run $(MAKE) $(RELEASE_MAKE_TARGETS)
RELEASE_MAKE_TARGETS +=

ifeq ($(SEMAPHORE), true)
ifeq ($(SEMAPHORE_PROJECT_ID),)
# The SEMAPHORE_PROJECT_ID variable is only set in sem2 environments
SEMAPHORE_2 := false
else
SEMAPHORE_2 := true
endif
endif

GIT_ROOT ?= $(CURDIR)
ifeq ($(SEMAPHORE_2),true)
# TODO: try to clean up .semaphore/semaphore.yml files.
# export GO111MODULE=on
# export "GOPATH=$(go env GOPATH)"
# export "SEMAPHORE_GIT_DIR=${GOPATH}/src/github.com/confluentinc/${SEMAPHORE_PROJECT_NAME}"
# export "PATH=${GOPATH}/bin:${PATH}:/usr/local/kubebuilder/bin:/usr/local/kubebuilder"
# mkdir -vp "${SEMAPHORE_GIT_DIR}" "${GOPATH}/bin"
# export SEMAPHORE_CACHE_DIR=/home/semaphore
ifeq ($(abspath $(SEMAPHORE_GIT_DIR)),$(SEMAPHORE_GIT_DIR))
GIT_ROOT := $(SEMAPHORE_GIT_DIR)
else
GIT_ROOT := $(HOME)/$(SEMAPHORE_GIT_DIR)
endif
# Place ci-bin inside the project as Semaphore 2 only allows caching resources within the project workspace.
# This needs to be different from $(GO_OUTDIR) so it doesn't get cleaned up by clean-go target.
CI_BIN := $(GIT_ROOT)/ci-bin
else ifeq ($(SEMAPHORE),true)
GIT_ROOT := $(SEMAPHORE_PROJECT_DIR)
CI_BIN := $(SEMAPHORE_CACHE_DIR)/bin
else ifeq ($(BUILDKITE),true)
CI_BIN := /tmp/bin
endif

# Where test reports get generated, used by testbreak reporting.
ifeq ($(SEMAPHORE),true)
BUILD_DIR := $(GIT_ROOT)/build
else
BUILD_DIR := /tmp/build
endif
export BUILD_DIR

HOST_OS := $(shell uname | tr A-Z a-z)

ifeq ($(BIN_PATH),)
ifeq ($(CI),true)
BIN_PATH := $(CI_BIN)
else
ifeq ($(HOST_OS),darwin)
BIN_PATH ?= /usr/local/bin
else
ifneq ($(wildcard $(HOME)/.local/bin/.),)
BIN_PATH ?= $(HOME)/.local/bin
else
BIN_PATH ?= $(HOME)/bin
endif
endif
endif
endif

XARGS := xargs
ifeq ($(HOST_OS),linux)
XARGS += --no-run-if-empty
endif

ifeq ($(CI),true)
# downstream things (like cpd CI) assume BIN_PATH exists
$(shell mkdir -p $(BIN_PATH) 2>/dev/null)
PATH := $(BIN_PATH):$(PATH)
export PATH
endif

# Git stuff
BRANCH_NAME ?= $(shell git rev-parse --abbrev-ref HEAD || true)
# Set RELEASE_BRANCH if we're on master or vN.N.x
RELEASE_BRANCH := $(shell echo $(BRANCH_NAME) | grep -E '^($(MASTER_BRANCH)|v[0-9]+\.[0-9]+\.x(-[0-9]+\.[0-9]+\.[0-9]+(-ce)?-SNAPSHOT)?)$$')
# assume the remote name is origin by default
GIT_REMOTE_NAME ?= origin

# Makefile called
MAKEFILE_NAME ?= Makefile
MAKE_ARGS := -f $(MAKEFILE_NAME)

# Determine if we're on a hotfix branch
ifeq ($(RELEASE_BRANCH),$(MASTER_BRANCH))
HOTFIX := false
else
HOTFIX := true
endif

ifeq ($(CI),true)
_ := $(shell test -d $(CI_BIN) || mkdir -p $(CI_BIN))
export PATH = $(CI_BIN):$(shell printenv PATH)
endif

.PHONY: update-mk-include
update-mk-include:
ifneq ($(shell git status --untracked-files=no --porcelain),)
	$(error git must be clean to update mk-include)
endif
	rm -rf mk-include
	git commit -a -m 'reset mk-include'
	git subtree add --prefix mk-include git@github.com:confluentinc/cc-mk-include.git master --squash

.PHONY: add-github-templates
add-github-templates:
	$(eval project_root := $(shell git rev-parse --show-toplevel))
	$(eval mk_include_relative_path := ../mk-include)
	$(if $(wildcard $(project_root)/.github/pull_request_template.md),$(a error ".github/pull_request_template.md already exists, try deleting it"),)
	$(if $(filter $(BRANCH_NAME),$(MASTER_BRANCH)),$(error "You must run this command from a branch: 'git checkout -b add-github-pr-template'"),)
	
	@mkdir -p $(project_root)/.github
	@ln -s $(mk_include_relative_path)/.github/pull_request_template.md $(project_root)/.github
	@git add $(project_root)/.github/pull_request_template.md
	@git commit \
		-m "Add .github template for PRs $(CI_SKIP)" \
		-m "Adds the .github/pull_request_template.md as described in [1]" \
		-m "linking to the shared template in \`mk-include\`." \
		-m "" \
		-m "[1] https://github.com/confluentinc/cc-mk-include/pull/113"
	
	@git show
	@echo "Template added."
	@echo "Create PR with 'git push && git log --format=%B -n 1 | hub pull-request -F -'"

.PHONY: add-paas-github-templates
add-paas-github-templates:
	$(eval project_root := $(shell git rev-parse --show-toplevel))
	$(eval mk_include_relative_path := ../mk-include)
	$(if $(wildcard $(project_root)/.github/pull_request_template.md),$(a error ".github/pull_request_template.md already exists, try deleting it"),)
	$(if $(filter $(BRANCH_NAME),$(MASTER_BRANCH)),$(error "You must run this command from a branch: 'git checkout -b add-github-pr-template'"),)
	
	@mkdir -p $(project_root)/.github
	@ln -s $(mk_include_relative_path)/.github/paas_pull_request_template.md $(project_root)/.github/pull_request_template.md
	@git add $(project_root)/.github/pull_request_template.md
	@git commit \
		-m "Add .github template for PRs $(CI_SKIP)" \
		-m "Adds the .github/pull_request_template.md as described in [1]" \
		-m "linking to the shared template in \`mk-include\`." \
		-m "" \
		-m "[1] https://github.com/confluentinc/cc-mk-include/pull/113"
	
	@git show
	@echo "Template added."
	@echo "Create PR with 'git push && git log --format=%B -n 1 | hub pull-request -F -'"

.PHONY: bats
bats:
	find . -name *.bats -exec bats {} \;

$(HOME)/.netrc:
ifeq ($(CI),true)
	$(error .netrc missing, can't authenticate to GitHub)
else
	$(shell bash -c 'echo .netrc missing, prompting for user input >&2')
	$(shell bash -c 'echo Enter Github credentials, if you use 2 factor authentication generate a personal access token for the password: https://github.com/settings/tokens >&2')
	$(eval user := $(shell bash -c 'read -p "GitHub Username: " user; echo $$user'))
	$(eval pass := $(shell bash -c 'read -s -p "GitHub Password: " pass; echo $$pass'))
	@printf "machine github.com\n\tlogin $(user)\n\tpassword $(pass)\n\nmachine api.github.com\n\tlogin $(user)\n\tpassword $(pass)\n" > $(HOME)/.netrc
	@echo
endif

.netrc: $(HOME)/.netrc
	cp $(HOME)/.netrc .netrc

.ssh: $(HOME)/.ssh
	cp -R $(HOME)/.ssh .ssh


.aws: $(HOME)/.aws
	cp -R $(HOME)/.aws .aws

.gitconfig: $(HOME)/.gitconfig
	cp $(HOME)/.gitconfig .gitconfig

.gcloud: $(HOME)/.config/gcloud
	mkdir -p .config
	cp -r $(HOME)/.config/gcloud/ .config/gcloud
