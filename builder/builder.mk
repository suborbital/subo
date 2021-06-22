# all paths are relative to project root
img = suborbital/builder-as:$(shell cat ./builder/.image-ver)

builder/docker: subo/docker builder/docker/rust builder/docker/swift builder/docker/as

builder/docker/publish: subo/docker/publish builder/docker/rs/publish builder/docker/swift/publish builder/docker/as/publish

# AssemblyScript docker targets
builder/docker/as:
	docker build . -f builder/docker/assemblyscript/Dockerfile -t $(img)

builder/docker/as/publish:
	docker buildx build . -f builder/docker/assemblyscript/Dockerfile --platform linux/amd64,linux/arm64 -t $(img) --push

# Rust docker targets
builder/docker/rust:
	docker build . -f builder/docker/rust/Dockerfile -t $(img)

builder/docker/rust/publish:
	docker buildx build . -f builder/docker/rust/Dockerfile --platform linux/amd64,linux/arm64 -t $(img) --push

# Swift docker targets
builder/docker/swift:
	docker build . -f builder/docker/swift/Dockerfile -t $(img)

builder/docker/swift/publish:
	docker buildx build . -f builder/docker/swift/Dockerfile --platform linux/amd64,linux/arm64 -t $(img) --push

.PHONY: builder/docker builder/docker/publish builder/docker/as builder/docker/as/publish builder/docker/rust builder/docker/rust/publish builder/docker/swift builder/docker/swift/publish