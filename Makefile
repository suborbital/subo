subo:
	go install ./subo

builder/rs:
	@$(MAKE) --no-print-directory -C builders/rust $@

builder/rs/%:
	@$(MAKE) --no-print-directory -C builders/rust $@

builder/swift:
	@$(MAKE) --no-print-directory -C builders/swift $@

builder/swift/%:
	@$(MAKE) --no-print-directory -C builders/swift $@

.PHONY: subo