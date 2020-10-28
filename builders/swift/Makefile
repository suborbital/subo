
builder/swift/test/wasm:
	swiftc --version
	TOOLCHAIN_PATH=$(dirname $(swiftenv which swiftc))/../
	mkdir -p .build
	cp lib.swift .build/lib.swift
	cat run.swift >> .build/lib.swift
	swiftc -target wasm32-unknown-wasi -sdk $(TOOLCHAIN_PATH)share/wasi-sysroot .build/lib.swift -o runnable.wasm -Xlinker --allow-undefined -Xlinker --export=allocate -Xlinker --export=deallocate -Xlinker --export=run_e


.PHONY: builder/swift/test/wasm