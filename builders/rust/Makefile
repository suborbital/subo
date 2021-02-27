rs-img = suborbital/builder-rs:$(shell cat .rs-ver)

builder/rs: subo/docker
	docker build ../../ -f ./Dockerfile -t $(rs-img)

builder/rs/publish:
	docker buildx build ../../ -f ./Dockerfile --platform linux/amd64,linux/arm64 -t $(rs-img) --push

builder/rs/test: builder/rs
	docker run -it --mount type=bind,source=$(PWD)/test-project/helloworld,target=/root/runnable $(rs-img)
	
.PHONY: builder/rs builder/rs/publish builder/rs/test builder/rs/cargo/publish subo/docker
