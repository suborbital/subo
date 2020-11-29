package context

import "fmt"

var nativeCommandsForLang = map[string][]string{
	"rust": {
		"cargo build --target wasm32-wasi --lib --release",
		"cp target/wasm32-wasi/release/{{ .UnderscoreName }}.wasm .",
	},
	"swift": {
		"xcrun --toolchain swiftwasm swift build --triple wasm32-unknown-wasi -Xlinker --allow-undefined -Xlinker --export=allocate -Xlinker --export=deallocate -Xlinker --export=run_e",
		"cp .build/debug/{{ .Name }}.wasm .",
	},
}

// NativeBuildCommands returns the native build commands needed to build a Runnable of a particular language
func NativeBuildCommands(lang string) ([]string, error) {
	cmds, exists := nativeCommandsForLang[lang]
	if !exists {
		return nil, fmt.Errorf("unable to build %s Runnables natively", lang)
	}

	return cmds, nil
}
