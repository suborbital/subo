as-img = suborbital/builder-as:$(shell cat ../.builder-ver)

builder/as: subo/docker
	docker build ../../ -f ./Dockerfile -t $(as-img)

builder/as/publish:
	docker buildx build ../../ -f ./Dockerfile --platform linux/amd64,linux/arm64 -t $(as-img) --push

builder/as/test: builder/as
	docker run -it --mount type=bind,source=$(PWD)/test-project/helloworld,target=/root/runnable $(as-img)
	
.PHONY: builder/as builder/as/publish builder/as/test builder/as/cargo/publish subo/docker
