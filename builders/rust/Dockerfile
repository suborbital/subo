FROM subo:dev as subo

FROM rust:1.46.0-slim-buster

WORKDIR /root
RUN mkdir runnable; mkdir suborbital

RUN rustup target install wasm32-wasi

COPY --from=subo /go/bin/subo /usr/local/bin
COPY ./suborbital /root/suborbital

ENV SUBO_DOCKER=1

WORKDIR /root/suborbital
RUN cargo build --target wasm32-wasi --lib --release

WORKDIR /root/runnable

ENTRYPOINT subo build --native .