# all paths are relative to project root
ver = $(shell cat ./builder/.image-ver)
tinygo_ver = $(shell cat ./builder/docker/tinygo/.tinygo-ver)

builder/docker: subo/docker builder/docker/rust builder/docker/swift builder/docker/as builder/docker/tinygo builder/docker/grain builder/docker/javascript

builder/docker/publish: subo/docker/publish builder/docker/rust/publish builder/docker/swift/publish builder/docker/as/publish builder/docker/tinygo/publish builder/docker/grain/publish builder/docker/javascript/publish

builder/docker/dev/publish: subo/docker/publish builder/docker/rust/dev/publish builder/docker/swift/dev/publish builder/docker/as/dev/publish builder/docker/tinygo/dev/publish builder/docker/grain/dev/publish builder/docker/javascript/dev/publish

# AssemblyScript docker targets
builder/docker/as:
	DOCKER_BUILDKIT=1 docker build . -f builder/docker/assemblyscript/Dockerfile -t suborbital/builder-as:$(ver)

builder/docker/as/publish:
	docker buildx build . -f builder/docker/assemblyscript/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-as:$(ver) --push

builder/docker/as/dev/publish:
	docker buildx build . -f builder/docker/assemblyscript/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-as:dev --push

# Rust docker targets
builder/docker/rust:
	DOCKER_BUILDKIT=1 docker build . -f builder/docker/rust/Dockerfile -t suborbital/builder-rs:$(ver)

builder/docker/rust/publish:
	docker buildx build . -f builder/docker/rust/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-rs:$(ver) --push

builder/docker/rust/dev/publish:
	docker buildx build . -f builder/docker/rust/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-rs:dev --push

# Swift docker targets
builder/docker/swift:
	DOCKER_BUILDKIT=1 docker build . -f builder/docker/swift/Dockerfile -t suborbital/builder-swift:$(ver)

builder/docker/swift/publish:
	docker buildx build . -f builder/docker/swift/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-swift:$(ver) --push

builder/docker/swift/dev/publish:
	docker buildx build . -f builder/docker/swift/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-swift:dev --push

# TinyGo (base) docker targets
builder/docker/tinygo-base:
	DOCKER_BUILDKIT=1 docker build . -f builder/docker/tinygo/Dockerfile.base -t suborbital/tinygo-base:$(tinygo_ver)

builder/docker/tinygo-base/publish:
	docker buildx build . -f builder/docker/tinygo/Dockerfile.base --platform linux/amd64,linux/arm64 -t suborbital/tinygo-base:$(tinygo_ver) --push

builder/docker/tinygo-base/dev/publish:
	docker buildx build . -f builder/docker/tinygo/Dockerfile.base --platform linux/amd64,linux/arm64 -t suborbital/tinygo-base:dev --push

# TinyGo (slim) docker targets
builder/docker/tinygo: builder/docker/tinygo-base
	DOCKER_BUILDKIT=1 docker build . -f builder/docker/tinygo/Dockerfile -t suborbital/builder-tinygo:$(ver)

builder/docker/tinygo/publish:
	docker buildx build . -f builder/docker/tinygo/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-tinygo:$(ver) --push

builder/docker/tinygo/dev/publish:
	docker buildx build . -f builder/docker/tinygo/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-tinygo:dev --push

# Grain docker targets
builder/docker/grain:
	docker buildx build . -f builder/docker/grain/Dockerfile --platform linux/amd64 -t suborbital/builder-gr:$(ver) --load

builder/docker/grain/publish:
	docker buildx build . -f builder/docker/grain/Dockerfile --platform linux/amd64 -t suborbital/builder-gr:$(ver) --push

builder/docker/grain/dev/publish:
	docker buildx build . -f builder/docker/grain/Dockerfile --platform linux/amd64 -t suborbital/builder-gr:dev --push

# JavaScript docker targets
builder/docker/javascript:
	DOCKER_BUILDKIT=1 docker build . -f builder/docker/javascript/Dockerfile -t suborbital/builder-js:$(ver)

builder/docker/javascript/publish:
	docker buildx build . -f builder/docker/javascript/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-js:$(ver) --push

builder/docker/javascript/dev/publish:
	docker buildx build . -f builder/docker/javascript/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-js:dev --push

.PHONY: builder/docker builder/docker/publish builder/docker/as builder/docker/as/publish builder/docker/rust builder/docker/rust/publish builder/docker/swift builder/docker/swift/publish builder/docker/tinygo builder/docker/tinygo/publish builder/docker/grain builder/docker/grain/publish builder/docker/javascript builder/docker/javascript/publish
