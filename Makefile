subo:
	go install ./subo

subo/dev:
	go install -tags=development ./subo

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

builder/as:
	@$(MAKE) --no-print-directory -C builders/assemblyscript $@

builder/as/%:
	@$(MAKE) --no-print-directory -C builders/assemblyscript $@

builders/publish: builder/rs/publish builder/swift/publish builder/as/publish

.PHONY: subo subo/docker