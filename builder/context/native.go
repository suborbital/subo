package context

import (
	"fmt"
	"runtime"
)

var nativeCommandsForLang = map[string]map[string][]string{
	"darwin": {
		"rust": {
			"cargo build --target wasm32-wasi --lib --release",
			"cp target/wasm32-wasi/release/{{ .UnderscoreName }}.wasm ./{{ .Name }}.wasm",
		},
		"swift": {
			"xcrun --toolchain swiftwasm swift build --triple wasm32-unknown-wasi -Xlinker --allow-undefined -Xlinker --export=allocate -Xlinker --export=deallocate -Xlinker --export=run_e -Xlinker --export=init",
			"cp .build/debug/{{ .Name }}.wasm .",
		},
		"assemblyscript": {
			"npm run asbuild",
		},
		"tinygo": {
			"tinygo build -o {{ .Name }}.wasm -target wasi .",
		},
	},
	"linux": {
		"rust": {
			"cargo build --target wasm32-wasi --lib --release",
			"cp target/wasm32-wasi/release/{{ .UnderscoreName }}.wasm ./{{ .Name }}.wasm",
		},
		"swift": {
			"swift build --triple wasm32-unknown-wasi -Xlinker --allow-undefined -Xlinker --export=allocate -Xlinker --export=deallocate -Xlinker --export=run_e -Xlinker --export=init",
			"cp .build/debug/{{ .Name }}.wasm .",
		},
		"assemblyscript": {
			"chmod -R +r ./",
			"npm run asbuild",
		},
		"tinygo": {
			"tinygo build -o {{ .Name }}.wasm -target wasi .",
		},
	},
}

// NativeBuildCommands returns the native build commands needed to build a Runnable of a particular language
func NativeBuildCommands(lang string) ([]string, error) {
	os := runtime.GOOS

	cmds, exists := nativeCommandsForLang[os][lang]
	if !exists {
		return nil, fmt.Errorf("unable to build %s Runnables natively", lang)
	}

	return cmds, nil
}
