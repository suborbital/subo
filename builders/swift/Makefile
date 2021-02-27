swift-img = suborbital/builder-swift:$(shell cat .swift-ver)

builder/swift:
	docker build ../../ -f ./Dockerfile -t $(swift-img)

builder/swift/publish: builder/swift
	docker push $(swift-img)

builder/swift/test: builder/swift
	docker run -it --mount type=bind,source=$(PWD)/test-project/swift-example,target=/root/runnable $(swift-img)

.PHONY: builder/swift builder/swift/publish builder/swift/test