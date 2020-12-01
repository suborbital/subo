rs-img = suborbital/builder-rs:$(shell cat .rs-ver)

builder/rs: subo/docker
	docker build . -t $(rs-img)

builder/rs/publish: builder/rs
	docker push $(rs-img)

builder/rs/test: builder/rs
	docker run -it --mount type=bind,source=$(PWD)/test-project/helloworld,target=/root/runnable $(rs-img)

builder/rs/cargo/publish:
	cd suborbital
	cargo publish --target wasm32-wasi

subo/docker:
	docker build ../../ -t subo:dev

.PHONY: builder/rs builder/rs/publish builder/rs/test builder/rs/cargo/publish subo/docker