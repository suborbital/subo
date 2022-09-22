include ./builder/builder.mk
include ./subo/release/release.mk

GO_INSTALL=go install -ldflags $(RELEASE_FLAGS)

subo:
	$(GO_INSTALL)

subo/dev:
	$(GO_INSTALL) -tags=development

subo/docker-bin:
	$(GO_INSTALL) -tags=docker

subo/docker:
	DOCKER_BUILDKIT=1 docker build . -t suborbital/subo:dev

subo/docker/publish:
	docker buildx build . --platform linux/amd64,linux/arm64 -t suborbital/subo:dev --push

subo/smoketest: subo
	./scripts/smoketest.sh

subo/toolchaintest: subo subo/docker
	./scripts/toolchaintest.sh

mod/replace/atmo:
	go mod edit -replace github.com/suborbital/atmo=$(HOME)/Workspaces/suborbital/atmo

tidy:
	go mod tidy && go mod download

lint:
	golangci-lint run ./...

test:
	go test ./...

.PHONY: subo subo/dev subo/docker-bin subo/docker subo/docker/publish subo/smoketest mod/replace/atmo tidy lint test
