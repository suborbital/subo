# Subo, the Suborbital CLI

Subo is the command-line helper for working with the Suborbital Development Platform. Subo is used to build Wasm Runnables, generate new projects and config files, and more over time.

**You do not need to install language-specific tools to get started with WebAssembly and Subo!** A Docker toolchain is supported (see below) that can build your Runnables without needing to install language toolchains.

## Installing
If you're on Mac (M1 or Intel), the easiest way to install is via `brew`:
```
brew tap suborbital/subo
brew install subo
```

On Intel Macs or Linux, you can also install `subo` using cURL (uses [gobinaries](https://gobinaries.com)):
```
curl -Ls https://subo.suborbital.dev | sh
```

## Verify installation
Verify subo was installed:
```
subo --help
```

Subo does not have official support for Windows.

## Alternative: install from source (requires Go)
To build and install `subo`, clone this repo and run:
```
make subo
```

## Getting started
**To get started with Subo, visit the [Get started guide](./docs/get-started.md).**

## Builders
This repo contains builders for the various languages supported by Wasm Runnables. A builder is a Docker image that can build Runnables into Wasm modules, and is used internally by `subo` to build your code! See the [builders](./builder/docker) directory for more.

## Platforms
The `subo` tool supports the following platforms and operating systems:
|  | x86_64 | arm64
| --- | --- | --- |
| Mac | ✅ | ✅ |
| Linux | ✅ | ✅ |
| Windows | 🚫 | 🚫 |
 
The language toolchains used by `subo` support the following platforms:
| | x86_64 | arm64 | Docker |
| --- | --- | --- | --- |
| Rust | ✅ | ✅ | ✅ |
| AssemblyScript | ✅ | ✅ | ✅ |
| Swift | ✅ | 🚫 | 🟡 (no arm64) |

## Contributing

Please read the [contributing guide](./CONTRIBUTING.md) to learn about how you can contribute to Subo! We welcome all types of contribution.

By the way, Subo is also the name of our mascot, and it's pronounced Sooooobo.

![SOS-Space_Panda-Dark-small](https://user-images.githubusercontent.com/5942370/129103528-8b013445-a8a2-44bb-8b39-65d912a66767.png)

Copyright Suborbital contributors 2021.
