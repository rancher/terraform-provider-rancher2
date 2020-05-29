RELEASE_POSTCOMMIT += push-docker-ubi
BUILD_TARGETS += build-docker-ubi

# if UBI_VERSION not defined suffix with tag ubi8 will be appended
ifndef UBI_VERSION
RHEL_UBI_TAG_NAME := $(BUILD_TAG)-ubi8
else
RHEL_UBI_TAG_NAME := $(BUILD_TAG)-$(UBI_VERSION)
endif
# RedHat release version scheme based on RedHat ubi
minor_version := $(subst .,$(_space),$(VERSION_NO_V))
RHEL_UBI_RELEASE_NUMBER := $(shell expr $(word 2,$(minor_version)))

.PHONY: build-docker-ubi
build-docker-ubi:
	cp $(HOME)/.netrc .netrc
	cp -R $(HOME)/.ssh .ssh
	docker build -f Dockerfile.ubi --no-cache --build-arg version=$(IMAGE_VERSION) --build-arg release=$(RHEL_UBI_RELEASE_NUMBER) -t $(RHEL_UBI_TAG_NAME) .
	rm -rf .netrc .ssh

.PHONY: tag-docker-ubi
tag-docker-ubi:
	@echo 'create docker tag $(IMAGE_VERSION)'
	docker tag $(RHEL_UBI_TAG_NAME) $(DOCKER_REPO)/$(RHEL_UBI_TAG_NAME)

.PHONY: push-docker-ubi
push-docker-ubi: tag-docker-ubi
	@echo 'push $(IMAGE_VERSION) to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(RHEL_UBI_TAG_NAME)
