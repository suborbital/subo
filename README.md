# Subo, the Suborbital CLI

Subo is the command-line helper for working with the Suborbital Development Platform. Subo is used to build Wasm Runnables, generate new projects and config files, and more over time.

**You do not need to install language-specific tools to get started with WebAssembly and Subo!** A Docker toolchain is supported (see below) that can build your Runnables without needing to install language toolchains.

## Installing
To install `subo`, clone this repo and run `make subo`. A version of Go that supports Modules is required. Package manager installations will be available soon.

You can also install with cURL (uses [gobinaries](https://gobinaries.com), does not support Apple Silicon):
```
curl -Ls https://subo.suborbital.dev | sh
```

Verify subo was installed:
```
subo --help
```

## Getting started
**To get started with Subo, visit the [Get started guide](./docs/get-started.md).**

## Builders
This repo contains builders for the various languages supported by Wasm Runnables. A builder is a Docker image that can build Runnables into Wasm modules, and is used internally by `subo` to build your code! See the [builders](./builders/) directory for more.

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
| Swift | ✅ | 🚫 | 🟡 (no arm64) |

## Contributing

Please read the [contributing guide](./CONTRIBUTING.md) to learn about how you can contribute to Subo! We welcome all types of contribution.

By the way, Subo is (in spirit) a chubby astronaut panda bear (with a retro Mercury-era vibe), and if any designer out there wants to illustrate them, the Suborbital contributors will find some way to compensate you for your time and effort. Also, it's pronounced Sooooobo.

Copyright Suborbital contributors 2021.
