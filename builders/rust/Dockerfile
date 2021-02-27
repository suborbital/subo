FROM subo:dev as subo

FROM rust:1.50.0-slim-buster

WORKDIR /root
RUN mkdir runnable; mkdir suborbital

RUN rustup target install wasm32-wasi
RUN true | cargo install lazy_static

COPY --from=subo /go/bin/subo /usr/local/bin

ENV SUBO_DOCKER=1

WORKDIR /root/runnable

ENTRYPOINT subo build --native .
