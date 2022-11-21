# Get started

Subo includes the WebAssembly toolchain for Suborbital projects.

The Suborbital Development Platform aims for Wasm to be a first-class citizen. `subo` is the toolchain for building Wasm modules for [E2Core](https://github.com/suborbital/e2core). The `subo` CLI can build Wasm modules, and can package several Wasm modules into a deployable bundle.

Building a modules in languages other than Go is designed to be simple and powerful:

```rust
impl runnable::Runnable for Example {
    fn run(&self, input: Vec<u8>) -> Option<Vec<u8>> {
        let in_string = String::from_utf8(input).unwrap();

        Some(String::from(format!("hello {}", in_string)).as_bytes().to_vec())
    }
}
```

subo will package your module into a Wasm module that can be used by E2Core and run just like any other module! You can see examples of modules in the [E2Core repository](https://github.com/suborbital/e2core/tree/main/sat/engine/testdata).

## Create a project

To create a new project for E2Core, use `subo create project <name>`. This will create a new folder which contains a Directive.yaml and an example module.

Full options for `create project`:

```console
Usage:
  subo create project <name> [flags]

Flags:
      --branch string        git branch to download templates from (default "main")
      --environment string   project environment name (your company's reverse domain (default "com.suborbital")
  -h, --help                 help for project
      --update-templates     update with the newest templates
```

## Create a module

To create a new module, use the create module command:

```console
> subo create module <name>
```

Rust is chosen by default, but if you prefer Swift, just pass `--lang=swift`! You can now use the module API to build your function. A directory is created for each module, and each contains a `.module.yaml` file that includes some metadata.

The full options for `create module`:

```console
Usage:
  subo create module <name> [flags]

Flags:
      --branch string      git branch to download templates from (default "main")
      --dir string         the directory to put the new module in (default "~/subo")
  -h, --help               help for module
      --lang string        the language of the new module (default "rust")
      --namespace string   the namespace for the new module (default "default")
      --repo string        git repo to download templates from (default "suborbital/templates")
      --update-templates   update with the newest module templates
```

## Building Wasm modules

**It is recommended that Docker be installed to build Wasm modules. See below if you do not have Docker installed.**

To build your module into a Wasm module for E2Core, use the build command:

```console
> subo build .
```

If the current working directory is a module, subo will build it. If the current directory contains many modules, subo will build them all. Any directory with a `.module.yaml` file is considered a module and will be built. Building modules is not fully tested on Windows.

## Bundles

By default, subo will write all of the modules in the current directory into a bundle. E2Core uses modules to help you build powerful web services by composing modules declaratively. If you want to skip bundling, you can pass `--no-bundle` to `subo build`

The resulting bundle can also be used with a Reactr instance by calling `h.HandleBundle({path/to/bundle})`. See the [Reactr Wasm instructions](https://github.com/suborbital/reactr/blob/master/docs/wasm.md) for details.

The full options for `build`:

```console
Usage:
  subo build [dir] [flags]

Flags:
      --builder-tag string   use the provided tag for builder images
      --docker               build your project's Dockerfile. It will be tagged {identifier}:{appVersion}
  -h, --help                 help for build
      --langs strings        build only modules for the listed languages (comma-seperated)
      --make string          execute the provided Make target before building the project bundle
      --mountpath string     if passed, the Docker builders will mount their volumes at the provided path
      --native               use native (locally installed) toolchain rather than Docker
      --no-bundle            if passed, a .wasm.zip bundle will not be generated
      --relpath subo build   if passed, the Docker builders will run subo build using the provided path, relative to '--mountpath'
```

## Building without Docker

If you prefer not to use Docker, you can use the `--native` flag. This will cause subo to use your local machine's toolchain to build modules instead of Docker containers. You will need to install the toolchains yourself:

- Rust: Install the latest Rust toolchain and the additional `wasm32-wasi` target.
- Swift: Install the [SwiftWasm](https://book.swiftwasm.org/getting-started/setup.html) toolchain. If using macOS, ensure XCode developer tools are installed (xcrun is required).

`subo` is continually evolving alongside [E2Core](https://github.com/suborbital/e2core).
