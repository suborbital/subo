# Subo, the Suborbital CLI

Subo is the command-line helper for working with the Suborbital Development Platform. Subo is used to build Wasm Runnables, generate new projects and config files, and more over time.

## Installing
To install `subo`, clone this repo and run `make subo`. A version of Go that supports Modules is required. Package manager installations will be available soon.

You can also install using [gobinaries](https://gobinaries.com/):
```
curl -sf https://gobinaries.com/suborbital/subo/subo | sh
```

## Getting started
**To get started with Subo, visit the [Get started guide](./docs/get-started.md).**

## Builders
This repo contains builders for the various languages supported by Wasm Runnables. A builder is a Docker image that can build Runnables into Wasm modules, and is used internally by `subo` to build your code! See the [builders](./builders/) directory for more.

By the way, Subo is (in spirit) a chubby astronaut panda bear (with a retro Mercury-era vibe), and if any designer out there wants to illustrate them, the Suborbital contributors will find some way to compensate you for your time and effort. Also, it's pronounced Sooooobo.

Copyright Suborbital contributors 2020.