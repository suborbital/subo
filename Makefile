subo:
	go install ./subo

subo/docker:
	docker build . -t subo:dev

builder/rs:
	@$(MAKE) --no-print-directory -C builders/rust $@

builder/rs/%:
	@$(MAKE) --no-print-directory -C builders/rust $@

builder/swift:
	@$(MAKE) --no-print-directory -C builders/swift $@

builder/swift/%:
	@$(MAKE) --no-print-directory -C builders/swift $@

.PHONY: subo