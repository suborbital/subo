rs-img = suborbital/builder-rs:$(shell cat .rs-ver)

builder/rs: subo/docker
	docker build . -t $(rs-img)

builder/rs/publish: subo/docker 
	docker buildx build . --platform linux/amd64,linux/arm64 -t $(rs-img)	

builder/rs/test: builder/rs
	docker run -it --mount type=bind,source=$(PWD)/test-project/helloworld,target=/root/runnable $(rs-img)

subo/docker:
	docker buildx build ../../ --platform linux/amd64,linux/arm64 -t subo:dev --load --no-cache

.PHONY: builder/rs builder/rs/publish builder/rs/test builder/rs/cargo/publish subo/docker
