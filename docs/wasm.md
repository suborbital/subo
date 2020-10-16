# Hive ❤️ Wasm

Subo includes the Wasm toolchain for [Hive](https://github.com/suborbital/hive).

The Suborbital Development Platform aims for Wasm to be a first-class citizen. `subo` is the toolchain for building Wasm Runnables for [Hive](https://github.com/suborbital/hive). The `subo` CLI can build Wasm Runnables, and can package several Wasm Runnables into a deployable bundle. It will soon be able to act as an all-in-one Wasm server, using Hive's FaaS functionality.

Writing a Runnable for Hive in languages other than Go is designed to be just as simple and powerful:
```rust
#[no_mangle]
pub fn run(input: Vec<u8>) -> Option<Vec<u8>> {
    let in_string = String::from_utf8(input).unwrap();

    Some(String::from(format!("hello {}", in_string)).as_bytes().to_vec())
}
```
subo will package your Runnable into a Wasm module that can be loaded into a Hive instance and run just like any other Runnable!

## Building Wasm Runnables
**Docker must be installed to build Wasm Runnables.**
The subo CLI builds your Runnable code into a Wasm module that can be loaded by Hive.

To build a Rust-based Runnable, see [helloworld-rs](./examples/helloworld-rs/README.md)

## Bundles
To build all of the Runnables in the current directory and bundle them all into a single `.wasm.zip` file, run `subo build --bundle`. The resulting bundle can be used with a Hive instance by calling `h.HandleBundle({path/to/bundle})`. See the [hive Wasm instructions](https://github.com/suborbital/hive/blob/master/Wasm.md) for details.

`subo` is under active development alongside [Hive](https://github.com/suborbital/hive) itself.

Copyright Suborbital contributors 2020

## FFI Runnable API
Hive provides an API which allows for communication between Wasm runnables and Hive. Full documentation is coming soon. This API currently has:
- The ability to make HTTP requests from Wasm Runnables (soon with built-in access controls to restrict network activity)

This API will soon have:
- The ability to read files from the host machine (with build-in access control)
- The ability to schedule new Hive jobs and get their results (similar to the Go Runnable API)