include ./builder/builder.mk
include ./cli/release/release.mk

GO_INSTALL=go install -ldflags $(RELEASE_FLAGS)

velo:
	$(GO_INSTALL)

velo/dev:
	$(GO_INSTALL) -tags=development

velo/docker-bin:
	$(GO_INSTALL) -tags=docker

velo/docker:
	DOCKER_BUILDKIT=1 docker build . -t suborbital/velo:dev

velo/docker/publish:
	docker buildx build . --platform linux/amd64,linux/arm64 -t suborbital/velo:dev --push

velo/smoketest: velo
	./scripts/smoketest.sh

mod/replace/atmo:
	go mod edit -replace github.com/suborbital/atmo=$(HOME)/Workspaces/suborbital/atmo

tidy:
	go mod tidy && go mod download

lint:
	golangci-lint run ./...

test:
	go test ./...

.PHONY: velo velo/dev velo/docker-bin velo/docker velo/docker/publish velo/smoketest mod/replace/atmo tidy lint test
