# all paths are relative to project root
include ./builder/docker/base/.env
image_ver = $(shell git describe --tags --dirty)

builder/docker/clean:
	rm ./builder/docker/base/.version

builder/docker: subo/docker builder/docker/rust builder/docker/swift builder/docker/as builder/docker/tinygo builder/docker/grain builder/docker/javascript

builder/docker/publish: subo/docker/publish builder/docker/rust/publish builder/docker/swift/publish builder/docker/as/publish builder/docker/tinygo/publish builder/docker/grain/publish builder/docker/javascript/publish

builder/docker/dev/publish: subo/docker/publish builder/docker/rust/dev/publish builder/docker/swift/dev/publish builder/docker/as/dev/publish builder/docker/tinygo/dev/publish builder/docker/grain/dev/publish builder/docker/javascript/dev/publish

# Force rebuild base image if .env or install.sh have changed
builder/docker/base/.version: ./builder/docker/base/.env ./builder/docker/base/install.sh ./builder/docker/base/Dockerfile
	docker build ./builder/docker/base\
		--no-cache\
		--build-arg BUILDER_VERSION=$(image_ver)\
		-t suborbital/build-pack:$(image_ver)

# Only rebuild docker if one of it's inputs changed
builder/docker/base: ./builder/docker/base/.version subo/docker
	@cat ./builder/docker/base/.env > ./builder/docker/base/.version

# AssemblyScript docker targets
builder/docker/as:
	docker build . -f builder/docker/assemblyscript/Dockerfile -t suborbital/builder-as:v${BUILDER_VERSION}

builder/docker/as/publish:
	docker buildx build . -f builder/docker/assemblyscript/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-as:v${BUILDER_VERSION} --push

builder/docker/as/dev/publish:
	docker buildx build . -f builder/docker/assemblyscript/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-as:dev --push

# Rust docker targets
builder/docker/rust: builder/docker/base
	docker build . -f builder/docker/rust/Dockerfile\
		--build-arg BUILDER_VERSION=$(image_ver)\
    	-t suborbital/builder-rs:$(image_ver)

builder/docker/rust/publish:
	docker buildx build . -f builder/docker/rust/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-rs:v${BUILDER_VERSION} --push

builder/docker/rust/dev/publish:
	docker buildx build . -f builder/docker/rust/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-rs:dev --push

# Swift docker targets
builder/docker/swift:
	docker build . -f builder/docker/swift/Dockerfile -t suborbital/builder-swift:v${BUILDER_VERSION}

builder/docker/swift/publish:
	docker buildx build . -f builder/docker/swift/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-swift:v${BUILDER_VERSION} --push

builder/docker/swift/dev/publish:
	docker buildx build . -f builder/docker/swift/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-swift:dev --push

builder/docker/tinygo-base/dev/publish:
	docker buildx build . -f builder/docker/tinygo/Dockerfile.base --platform linux/amd64,linux/arm64 -t suborbital/tinygo-base:dev --push

# TinyGo (slim) docker targets
builder/docker/tinygo: builder/docker/base
	docker build . -f builder/docker/tinygo/Dockerfile\
 		--build-arg BUILDER_VERSION=$(image_ver)\
 		-t suborbital/builder-tinygo:$(image_ver)

builder/docker/tinygo/publish:
	docker buildx build . -f $(CURDIR)/docker/tinygo/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-tinygo:v${BUILDER_VERSION} --push

builder/docker/tinygo/dev/publish:
	docker buildx build . -f builder/docker/tinygo/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-tinygo:dev --push

# Grain docker targets
builder/docker/grain:
	docker buildx build . -f builder/docker/grain/Dockerfile --platform linux/amd64 -t suborbital/builder-gr:v${BUILDER_VERSION} --load

builder/docker/grain/publish:
	docker buildx build . -f builder/docker/grain/Dockerfile --platform linux/amd64 -t suborbital/builder-gr:v${BUILDER_VERSION} --push

builder/docker/grain/dev/publish:
	docker buildx build . -f builder/docker/grain/Dockerfile --platform linux/amd64 -t suborbital/builder-gr:dev --push

# JavaScript docker targets
builder/docker/javy: builder/docker/rust
	docker build . -f builder/docker/javy/Dockerfile\
 		--build-arg BUILDER_VERSION=${BUILDER_VERSION}\
 		-t suborbital/javy:$(image_ver)

builder/docker/javy/publish:
	docker buildx build . -f builder/docker/javy/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/javy:v${BUILDER_VERSION} --push

builder/docker/javy/dev/publish:
	docker buildx build . -f builder/docker/javy/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/javy:dev --push

# JavaScript docker targets
builder/docker/javascript: builder/docker/javy
	docker build . -f builder/docker/javascript/Dockerfile\
 		--build-arg BUILDER_VERSION=${BUILDER_VERSION}\
 		-t suborbital/builder-js:$(image_ver)

builder/docker/javascript/publish:
	docker buildx build . -f builder/docker/javascript/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-js:v${BUILDER_VERSION} --push

builder/docker/javascript/dev/publish:
	docker buildx build . -f builder/docker/javascript/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-js:dev --push

.PHONY: builder/docker builder/docker/publish builder/docker/as builder/docker/as/publish builder/docker/base builder/docker/rust builder/docker/rust/publish builder/docker/swift builder/docker/swift/publish builder/docker/tinygo builder/docker/tinygo/publish builder/docker/grain builder/docker/grain/publish builder/docker/javascript builder/docker/javascript/publish
