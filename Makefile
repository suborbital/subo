include ./builder/builder.mk

subo:
	go install ./subo

subo/dev:
	go install -tags=development ./subo

subo/docker:
	docker buildx build . --platform linux/amd64 -t suborbital/subo:dev --load

subo/docker/publish:
	docker buildx build . --platform linux/amd64,linux/arm64 -t suborbital/subo:dev --push

mod/replace/atmo:
	go mod edit -replace github.com/suborbital/atmo=$(HOME)/Workspaces/suborbital/atmo

.PHONY: subo subo/docker