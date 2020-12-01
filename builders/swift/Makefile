swift-img = suborbital/builder-swift:$(shell cat .swift-ver)

builder/swift: subo/docker
	docker build . -t $(swift-img)

builder/swift/publish: builder/swift
	docker push $(swift-img)

builder/swift/test: builder/swift
	docker run -it --mount type=bind,source=$(PWD)/test-project/swift-example,target=/root/runnable $(swift-img)

builder/swift/cargo/publish:
	cd suborbital
	cargo publish --target wasm32-wasi

subo/docker:
	docker build ../../ -t subo:dev

.PHONY: builder/swift builder/swift/publish builder/swift/test builder/swift/cargo/publish