#! /bin/bash

set -e

cp ../rs-wasm/* ./
cp *.rs ./src

name=$(./target/debug/tomltool)
echo "building $name wasm package"

cargo build --target wasm32-wasi --lib --release

underscorename=$(echo $name | tr - _)

cp target/wasm32-wasi/release/${underscorename}.wasm ../rs-wasm/$name.wasm