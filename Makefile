subo:
	go install ./subo

builder/rs:
	@$(MAKE) --no-print-directory -C builders/rust $@

builder/rs/%:
	@$(MAKE) --no-print-directory -C builders/rust $@

.PHONY: subo