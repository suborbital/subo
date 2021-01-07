# Get started

Subo includes the WebAssembly toolchain for Suborbital projects.

The Suborbital Development Platform aims for Wasm to be a first-class citizen. `subo` is the toolchain for building Wasm Runnables for [Hive](https://github.com/suborbital/hive) and [Atmo](https://github.com/suborbital/atmo). The `subo` CLI can build Wasm Runnables, and can package several Wasm Runnables into a deployable bundle.

Building a Runnable in languages other than Go is designed to be simple and powerful:
```rust
impl runnable::Runnable for Example {
    fn run(&self, input: Vec<u8>) -> Option<Vec<u8>> {
        let in_string = String::from_utf8(input).unwrap();
    
        Some(String::from(format!("hello {}", in_string)).as_bytes().to_vec())
    }
}
```
subo will package your Runnable into a Wasm module that can be used by Hive or Atmo and run just like any other Runnable! You can see examples of Runnables in the [test project](../test-project).

## Create a project
To create a new project for Atmo or Hive, use `subo create project <name>`. This will create a new folder which contains a Directive.yaml and an example Runnable.

Full options for `create project`:
```
create a new project for Atmo or Hive

Usage:
  subo create project <name> [flags]

Flags:
      --branch string      git branch to download templates from (default "main")
  -h, --help               help for project
      --update-templates   update with the newest templates
```

## Create a Runnable
To create a new Runnable, use the create runnable command:
```
> subo create runnable <name>
```
Rust is chosen by default, but if you prefer Swift, just pass `--lang=swift`! You can now use the Runnable API to build your function. A directory is created for each Runnable, and each contains a `.runnable.yaml` file that includes some metadata.

The full options for `create runnable`:
```
Usage:
  subo create <name> [flags]

Flags:
      --branch string      git branch to download templates from (default "main")
      --dir string         the directory to put the new runnable in (default "/Users/cohix-16/Workspaces/suborbital/subo")
  -h, --help               help for create
      --lang string        the language of the new runnable (default "rust")
      --namespace string   the namespace for the new runnable (default "default")
      --update-templates   update with the newest runnable templates
```

## Building Wasm Runnables
**It is reccomended that Docker be installed to build Wasm Runnables. See below if you do not have Docker installed.**
 
To build your Runnable into a Wasm module for Hive or Atmo, use the build command:
```
> subo build .
```
If the current working directory is a Runnable, subo will build it. If the current directory contains many runnables, subo will build them all. Any directory with a `.runnable.yaml` file is considered a Runnable and will be built. Building Runnables is not fully tested on Windows.

## Bundles
To build all of the Runnables in the current directory and bundle them all into a single `.wasm.zip` file, run `subo build . --bundle`. Atmo uses Runnable bundles to help you build powerful web services by composing Runnables declaratively.

The resulting bundle can also be used with a Hive instance by calling `h.HandleBundle({path/to/bundle})`. See the [hive Wasm instructions](https://github.com/suborbital/hive/blob/master/Wasm.md) for details.

The full options for `build`:
```
Usage:
  subo build [dir] [flags]

Flags:
      --bundle   if passed, bundle all resulting runnables into a deployable .wasm.zip bundle
  -h, --help     help for build
      --native   if passed, build runnables using native toolchain rather than Docker
```

## Building without Docker
If you prefer not to use Docker, you can use the `--native` flag. This will cause subo to use your local machine's toolchain to build Runnables instead of Docker containers. You will need to install the toolchains yourself:
- Rust: Install the latest Rust toolchain and the additional `wasm32-wasi` target.
- Swift: Install the [SwiftWasm](https://book.swiftwasm.org/getting-started/setup.html) toolchain. If using macOS, ensure XCode developer tools are installed (xcrun is required).

`subo` is continually evolving alongside [Hive](https://github.com/suborbital/hive) and [Atmo](https://github.com/suborbital/atmo).

## Suborbital Runnable API
Hive provides an [API](https://github.com/suborbital/hive-wasm) which gives Wasm Runnables the ability to access resources and communicate with the host application. Full documentation is coming soon. This API currently has:
- The ability to make HTTP requests from Wasm Runnables (soon with built-in access controls to restrict network activity) (Rust)
- Logging abilities (Rust, Swift)
- Access to persistent cache (Rust, Swift)

This API will soon have:
- The ability to read static files packaged into Runnable bundles
- The ability to render templates
- Database access