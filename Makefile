OUTPUT_FOLDER ?= out
BINARY_NAME ?= nodeport-deleter
REPOSITORY ?= index.docker.io/kangcode/cilium-nodeport-deleter
TAG ?= latest

.PHONY: all
all: build package push

.PHONY: build
build:
	CGO_ENABLED=0 go build -o $(OUTPUT_FOLDER)/$(BINARY_NAME) ./
	chmod 555 $(OUTPUT_FOLDER)/$(BINARY_NAME)

.PHONY: package
package:
	docker build \
		--build-arg OUTPUT_FOLDER=$(OUTPUT_FOLDER) \
		--build-arg BINARY_NAME=$(BINARY_NAME) \
		-t $(REPOSITORY):$(TAG) \
		.

.PHONY: push
push:
	docker push $(REPOSITORY):$(TAG)
	docker push $(REPOSITORY):latest
