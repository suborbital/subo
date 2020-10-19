# This Dockerfile (and everything in rs-builder) is internal to the tool and doesn't need to be changed.

FROM rust:1.46.0-slim-buster

WORKDIR /root
RUN mkdir rs-builder; mkdir rs-wasm
WORKDIR /root/rs-builder

ADD . .

RUN rustup target install wasm32-wasi

RUN cargo build --bin tomltool
RUN cargo build --target wasm32-wasi --lib --release

RUN rm ./src/run.rs
RUN mv Empty.toml Cargo.toml

ENTRYPOINT ./build.sh