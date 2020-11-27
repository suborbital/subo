rs-img = suborbital/hive-rs:$(shell cat .rs-ver)

builder/rs:
	docker build . -t $(rs-img)

builder/rs/publish/docker: build/docker
	docker push $(rs-img)

builder/rs/example/wasm:
	cargo build --target wasm32-wasi --lib --release

builder/rs/test/docker: build/docker
	docker run -it --mount type=bind,source="$(PWD)"/helloworld-rs,target=/root/rs-wasm $(rs-img)

.PHONY: builder/rs builder/rs/publish/docker builder/rs/test/wasm builder/rs/test/docker