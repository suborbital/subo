subo:
	go install ./subo

subo/docker:
	docker build . -t subo:dev

subo/static:
	CGO_ENABLED=1 go install -ldflags "-linkmode external -extldflags -static" -a ./subo

builder/rs:
	@$(MAKE) --no-print-directory -C builders/rust $@

builder/rs/%:
	@$(MAKE) --no-print-directory -C builders/rust $@

builder/swift:
	@$(MAKE) --no-print-directory -C builders/swift $@

builder/swift/%:
	@$(MAKE) --no-print-directory -C builders/swift $@

.PHONY: subo