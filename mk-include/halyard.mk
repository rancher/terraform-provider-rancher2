# Addresses to halyard services
HALYARD_DEPLOYER_ADDRESS ?= halyard-deployer.prod.halyard.confluent.cloud:9090
HALYARD_RELEASE_ADDRESS ?= halyard-release.prod.halyard.confluent.cloud:9090
HALYARD_RENDERER_ADDRESS ?= halyard-renderer.prod.halyard.confluent.cloud:9090

# Determin which halyard services to auto bump source version
# List of halyard service files, default all.  All environments in these files will be bumped
HALYARD_SERVICE_FILES ?= $(wildcard .halyard/*.yaml)
# List of halyard service files with environments, defaults none.
# NOTE: This disables HALYARD_SERVICE_FILES, it's either full auto or full manual.
# NOTE: Apply always applies all files in HALYARD_SERVICE_FILES since it won't create new env
#       versions if there's nothing changed.
# Format: .halyard/service.yaml=env1 .halyard/service.yaml=env2 etc.
HALYARD_SERVICE_FILES_ENVS ?=
# Version to set source version to, defaults to current clean version without a v.
HALYARD_SOURCE_VERSION ?= $(BUMPED_CLEAN_VERSION)
# List of service/environments to automatically install on release, defaults none.
# Format: service=env service=env2 service2=env
HALYARD_INSTALL_SERVICE_ENVS ?=

# Only create a tmpdir on CI
ifeq ($(CI),true)
HAL_TMPDIR := $(shell mktemp -d 2>/dev/null || mktemp -d -t 'halyard')
else
HAL_TMPDIR := $(PWD)
endif

# setup halctl cmd
HALYARD_VERSION ?= latest
HALCTL_ARGS ?=
HALYARD_IMAGE ?= confluent-docker.jfrog.io/confluentinc/halyard:$(HALYARD_VERSION)
_halctl_opts := --deployer-address $(HALYARD_DEPLOYER_ADDRESS)
_halctl_opts += --release-address $(HALYARD_RELEASE_ADDRESS)
_halctl_opts += --renderer-address $(HALYARD_RENDERER_ADDRESS)
_halctl_opts += $(HALCTL_ARGS)
_halctl_docker_opts := --user $(shell id -u):$(shell id -g) --rm -t
_halctl_docker_opts += -v $(PWD):/work -v $(HOME)/.halctl:/.halctl -w /work
ifeq ($(CI),true)
_halctl_docker_opts += -v $(HAL_TMPDIR):$(HAL_TMPDIR)
endif
HALCTL ?= docker run $(_halctl_docker_opts) $(HALYARD_IMAGE) $(_halctl_opts)

INIT_CI_TARGETS += halyard-cache-image
RELEASE_PRECOMMIT += halyard-set-source-version
RELEASE_POSTCOMMIT += halyard-apply-services halyard-install-services

.PHONY: show-halyard
## Show Halyard Variables
show-halyard:
	@echo "HALYARD_SERVICE_FILES:        $(HALYARD_SERVICE_FILES)"
	@echo "HALYARD_SERVICE_FILES_ENVS:   $(HALYARD_SERVICE_FILES_ENVS)"
	@echo "HALYARD_INSTALL_SERVICE_ENVS: $(HALYARD_INSTALL_SERVICE_ENVS)"
	@echo "HALYARD_SOURCE_VERSION:       $(HALYARD_SOURCE_VERSION)"
	@echo "HALCTL:                       $(HALCTL)"

# target for caching the halyard docker image on semaphore
.PHONY: halyard-cache-image
halyard-cache-image:
	cache restore halyard-image
	test ! -f halyard-image.tgz || docker load -i halyard-image.tgz
	mv halyard-image.tgz halyard-image-prev.tgz || echo dummy > halyard-image-prev.tgz
	docker pull $(HALYARD_IMAGE)
	docker save $(HALYARD_IMAGE) | gzip --no-name > halyard-image.tgz
	cmp halyard-image-prev.tgz halyard-image.tgz || cache store halyard-image halyard-image.tgz
	rm -f halyard-image*.tgz

$(HOME)/.halctl:
	mkdir $(HOME)/.halctl

.PHONY: halctl
## Run halctl in the halyard docker image
halctl: $(HOME)/.halctl
	@$(HALCTL) $(HALCTL_ARGS)

.PHONY: halyard-set-source-version
ifeq ($(HALYARD_SERVICE_FILES_ENVS),)
halyard-set-source-version: $(HALYARD_SERVICE_FILES:%=set.%)
else
halyard-set-source-version: $(HALYARD_SERVICE_FILES_ENVS:%=set.%)
endif

.PHONY: $(HALYARD_SERVICE_FILES:%=set.%)
$(HALYARD_SERVICE_FILES:%=set.%): $(HOME)/.halctl
	$(HALCTL) release set-file-version -v $(HALYARD_SOURCE_VERSION) -f $(@:set.%=%)
	git add $(@:set.%=%)

.PHONY: $(HALYARD_SERVICE_FILES_ENVS:%=set.%)
$(HALYARD_SERVICE_FILES_ENVS:%=set.%): $(HOME)/.halctl
	@$(eval fpath := $(word 1,$(subst =, ,$@)))
	@$(eval env := $(word 2,$(subst =, ,$@)))
	$(HALCTL) release set-file-version -v $(HALYARD_SOURCE_VERSION) -f $(fpath) -e $(env)
	git add $(fpath)

.PHONY: halyard-apply-services
halyard-apply-services: $(HALYARD_SERVICE_FILES:%=apply.%)

.PHONY: $(HALYARD_SERVICE_FILES:%=apply.%)
$(HALYARD_SERVICE_FILES:%=apply.%): $(HOME)/.halctl
	$(HALCTL) release apply -f $(@:apply.%=%) --output-dir $(HAL_TMPDIR)

cc-releases:
	git clone git@github.com:confluentinc/cc-releases.git

.PHONY: update-cc-releases
update-cc-releases:
	git -C cc-releases checkout master
	git -C cc-releases pull

commit-cc-releases:
	git -C cc-releases diff --exit-code --cached --name-status || \
	git -C cc-releases commit -m "chore: auto update"
	rm -rf cc-releases

.PHONY: halyard-install-services
halyard-install-services: cc-releases update-cc-releases $(HALYARD_INSTALL_SERVICE_ENVS:%=install.%) commit-cc-releases

.PHONY: $(HALYARD_INSTALL_SERVICE_ENVS:%=install.%)
$(HALYARD_INSTALL_SERVICE_ENVS:%=install.%): $(HOME)/.halctl
	$(eval svc := $(word 1,$(subst =, ,$(@:install.%=%))))
	$(eval env := $(word 2,$(subst =, ,$(@:install.%=%))))
	$(eval ver := $(shell cat $(HAL_TMPDIR)/$(svc)/$(env)))
	$(HALCTL) release set-file-install-version -v $(ver) -f cc-releases/services/$(svc)/$(env).yaml
	git -C cc-releases add cc-releases/services/$(svc)/$(env).yaml
