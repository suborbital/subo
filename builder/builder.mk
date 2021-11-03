# all paths are relative to project root
ver = $(shell cat ./builder/.image-ver)

builder/docker: subo/docker builder/docker/rust builder/docker/swift builder/docker/as

builder/docker/publish: subo/docker/publish builder/docker/rust/publish builder/docker/swift/publish builder/docker/as/publish

builder/docker/dev/publish: subo/docker/publish builder/docker/rust/dev/publish builder/docker/swift/dev/publish builder/docker/as/dev/publish

# AssemblyScript docker targets
builder/docker/as:
	docker build . -f builder/docker/assemblyscript/Dockerfile -t suborbital/builder-as:$(ver)

builder/docker/as/publish:
	docker buildx build . -f builder/docker/assemblyscript/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-as:$(ver) --push

builder/docker/as/dev/publish:
	docker buildx build . -f builder/docker/assemblyscript/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-as:dev --push

# Rust docker targets
builder/docker/rust:
	docker build . -f builder/docker/rust/Dockerfile -t suborbital/builder-rs:$(ver)

builder/docker/rust/publish:
	docker buildx build . -f builder/docker/rust/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-rs:$(ver) --push

builder/docker/rust/dev/publish:
	docker buildx build . -f builder/docker/rust/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-rs:dev --push

# Swift docker targets
builder/docker/swift:
	docker build . -f builder/docker/swift/Dockerfile -t suborbital/builder-swift:$(ver)

builder/docker/swift/publish:
	docker buildx build . -f builder/docker/swift/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-swift:$(ver) --push

builder/docker/swift/dev/publish:
	docker buildx build . -f builder/docker/swift/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-swift:dev --push

# tinygo docker targets
builder/docker/tinygo:
	docker build . -f builder/docker/tinygo/Dockerfile -t suborbital/builder-tinygo:$(ver)

builder/docker/tinygo/publish:
	docker buildx build . -f builder/docker/tinygo/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-tinygo:$(ver) --push

builder/docker/tinygo/dev/publish:
	docker buildx build . -f builder/docker/tinygo/Dockerfile --platform linux/amd64,linux/arm64 -t suborbital/builder-tinygo:dev --push

# Grain docker targets
builder/docker/grain:
	docker buildx build . -f builder/docker/grain/Dockerfile --platform linux/amd64 -t suborbital/builder-gr:$(ver) --load

builder/docker/grain/publish:
	docker buildx build . -f builder/docker/grain/Dockerfile --platform linux/amd64 -t suborbital/builder-gr:$(ver) --push

.PHONY: builder/docker builder/docker/publish builder/docker/as builder/docker/as/publish builder/docker/rust builder/docker/rust/publish builder/docker/swift builder/docker/swift/publish builder/docker/tinygo builder/docker/tinygo/publish builder/docker/grain builder/docker/grain/publish
