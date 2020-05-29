_empty :=
_space := $(_empty) $(empty)

# Use this variable to specify a different make utility (e.g. remake --profile)
# Note: not using $(MAKE) here since that runs inside container (different OS)
DOCKER_MAKE ?= make

# List of base images, cannot have colons, replace with a bang
DOCKER_BASE_IMAGES ?= $(subst :,!,$(shell grep FROM Dockerfile | cut -d' ' -f2))

# Use this variable to specify docker build options
DOCKER_BUILD_OPTIONS ?=
ifeq ($(CI),true)
	DOCKER_BUILD_OPTIONS += --no-cache
endif

# Image Name
IMAGE_NAME ?= unknown
ifeq ($(IMAGE_NAME),unknown)
$(error IMAGE_NAME must be set)
endif

# Image Version
#  If we're on CI and a release branch, build with the bumped version
ifeq ($(CI),true)
ifneq ($(RELEASE_BRANCH),$(_empty))
IMAGE_VERSION ?= $(BUMPED_VERSION)
else
IMAGE_VERSION ?= $(VERSION)
endif
else
IMAGE_VERSION ?= $(VERSION)
endif

IMAGE_REPO ?= confluentinc
ifeq ($(IMAGE_REPO),$(_empty))
BUILD_TAG ?= $(IMAGE_NAME):$(IMAGE_VERSION)
BUILD_TAG_LATEST ?= $(IMAGE_NAME):latest
else
BUILD_TAG ?= $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_VERSION)
BUILD_TAG_LATEST ?= $(IMAGE_REPO)/$(IMAGE_NAME):latest
endif

DOCKER_REPO ?= confluent-docker.jfrog.io

# Set targets for standard commands
INIT_CI_TARGETS += cache-docker-base-images
RELEASE_POSTCOMMIT += push-docker
BUILD_TARGETS += build-docker
CLEAN_TARGETS += clean-images

DOCKER_BUILD_PRE ?=
DOCKER_BUILD_POST ?=

.PHONY: show-docker
## Show docker variables
show-docker:
	@echo "DOCKER_BASE_IMAGES: $(DOCKER_BASE_IMAGES)"
	@echo "IMAGE_NAME: $(IMAGE_NAME)"
	@echo "IMAGE_VERSION: $(IMAGE_VERSION)"
	@echo "IMAGE_REPO: $(IMAGE_REPO)"
	@echo "BUILD_TAG: $(BUILD_TAG)"
	@echo "BUILD_TAG_LATEST: $(BUILD_TAG_LATEST)"
	@echo "DOCKER_REPO: $(DOCKER_REPO)"

.PHONY: docker-login
## Login to docker Artifactory
docker-login:
ifeq ($(DOCKER_USER)$(DOCKER_APIKEY),$(_empty))
	@jq -e '.auths."confluent-docker.jfrog.io"' $(HOME)/.docker/config.json 2>&1 >/dev/null ||\
		(echo "confluent-docker.jfrog.io not logged in, Username and Password not found in environment, prompting for login:" && \
		 docker login confluent-docker.jfrog.io)
else
	@jq -e '.auths."confluent-docker.jfrog.io"' $(HOME)/.docker/config.json 2>&1 >/dev/null ||\
		docker login confluent-docker.jfrog.io --username $(DOCKER_USER) --password $(DOCKER_APIKEY)
endif

.PHONY: cache-docker-base-images $(DOCKER_BASE_IMAGES:%=docker-cache.%)
## On Semaphore, use the cache to store/restore docker image to reduce transfer costs.
## - use gzip --no-name so the bits are deterministic,
## - always pull, this checks for updates, e.g. 'latest' tag could have been updated,
## - update cache if bits are different.
cache-docker-base-images: docker-login $(DOCKER_BASE_IMAGES:%=docker-cache.%)
$(DOCKER_BASE_IMAGES:%=docker-cache.%):
	$(eval image := $(subst !,:,$(@:docker-cache.%=%)))
	cache restore $(image)
	test ! -f base-image.tgz || docker load -i base-image.tgz
	mv base-image.tgz base-image-prev.tgz || echo dummy > base-image-prev.tgz
	docker pull $(image)
	docker save $(image) | gzip --no-name > base-image.tgz
	cmp base-image-prev.tgz base-image.tgz || cache store $(image) base-image.tgz
	rm -f base-image*.tgz

.PHONY: build-docker
ifeq ($(BUILD_DOCKER_OVERRIDE),)
## Build just the docker image
build-docker: .gitconfig .netrc .ssh $(DOCKER_BUILD_PRE)
	docker build $(DOCKER_BUILD_OPTIONS) \
		--build-arg version=$(IMAGE_VERSION) \
		-t $(BUILD_TAG) .
	rm -rf .netrc .ssh
ifeq ($(CI),true)
	docker image save $(BUILD_TAG) | gzip | \
		artifact push project /dev/stdin -d docker/$(BRANCH_NAME)/$(IMAGE_VERSION).tgz --force
endif
ifneq ($(DOCKER_BUILD_POST),)
	$(MAKE) $(MAKE_ARGS) $(DOCKER_BUILD_POST)
endif
else
build-docker: $(BUILD_DOCKER_OVERRIDE)
endif

.PHONY: restore-docker-version
ifeq ($(RESTORE_DOCKER_OVERRIDE),)
restore-docker-version:
ifeq ($(CI),true)
	artifact pull project docker/$(BRANCH_NAME)/$(IMAGE_VERSION).tgz -d /dev/stdout --force | \
		gunzip | docker image load
endif
else
restore-docker-version: $(RESTORE_DOCKER_OVERRIDE)
endif

.PHONY: tag-docker
tag-docker: tag-docker-latest tag-docker-version

.PHONY: tag-docker-latest
tag-docker-latest:
	@echo 'create docker tag latest'
	docker tag $(BUILD_TAG) $(DOCKER_REPO)/$(BUILD_TAG_LATEST)

.PHONY: tag-docker-version
tag-docker-version:
	@echo 'create docker tag $(IMAGE_VERSION)'
	docker tag $(BUILD_TAG) $(DOCKER_REPO)/$(BUILD_TAG)

.PHONY: push-docker
ifeq ($(PUSH_DOCKER_OVERRIDE),)
push-docker: push-docker-version push-docker-latest
else
push-docker: $(PUSH_DOCKER_OVERRIDE)
endif

.PHONY: ci-docker
ci-docker: .netrc .ssh .aws .gitconfig .gcloud docker-pull-base $(DOCKER_BUILD_PRE)
	docker build -t $(IMAGE_NAME)-ci --file Dockerfile-ci .
		docker run --tty \
		--network host \
		--env DOCKER_APIKEY \
		--env DOCKER_USER \
		--env HELM_APIKEY \
		--env HELM_USER \
		--env DATABASE_POSTGRESQL_USERNAME \
		--env DATABASE_POSTGRESQL_PASSWORD \
		--env SEMAPHORE_PROJECT_NAME \
		--env BRANCH_NAME \
		--env SEMAPHORE_BUILD_NUMBER \
		--env SEMAPHORE_CACHE_DIR \
		--env TEST_REPORT_FILE=/root/TEST-result.xml \
		--env SEED_POSTGRES_URL=postgres://$(DATABASE_POSTGRESQL_USERNAME):$(DATABASE_POSTGRESQL_PASSWORD)@localhost:5432 \
		--env CODECOV=true \
		--env CI=true \
		--volume /var/run/docker.sock:/var/run/docker.sock \
		--volume /var/run/postgresql:/var/run/postgresql \
		--volume $(shell pwd):/root/$(MODULE_NAME) \
		--volume $(GOPATH)/pkg/mod/cache/download:/go/pkg/mod/cache/download \
		--volume $(SEMAPHORE_CACHE_DIR):/root/$(MODULE_NAME)/.semaphore-cache \
		--volume $(GOPATH)/pkg/mod/cache/download:/root/$(MODULE_NAME)/.gomodcache \
		$(shell bash <(curl -s https://codecov.io/env)) \
		$(IMAGE_NAME)-ci $(DOCKER_MAKE) -C /root/$(MODULE_NAME) init-ci deps test build release-ci

.PHONY: push-docker-latest
push-docker-latest: tag-docker-latest
	@echo 'push latest to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(BUILD_TAG_LATEST)

.PHONY: push-docker-version
## Push the current version of docker to artifactory
push-docker-version: restore-docker-version tag-docker-version
	@echo 'push $(IMAGE_VERSION) to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(BUILD_TAG)

.PHONY: clean-images
clean-images:
	docker images -q -f label=io.confluent.caas=true -f reference='*$(IMAGE_NAME)' | uniq | $(XARGS) docker rmi -f

.PHONY: clean-all
clean-all:
	docker images -q -f label=io.confluent.caas=true | uniq | $(XARGS) docker rmi -f
