FROM subo:dev as subo

FROM rust:1.46.0-slim-buster

WORKDIR /root
RUN mkdir runnable; mkdir suborbital

RUN rustup target install wasm32-wasi

COPY --from=subo /go/bin/subo /usr/local/bin
COPY --from=subo /go/pkg/mod/github.com/wasmerio/wasmer-go@v0.3.2-0.20200903143934-299b1d478cde/wasmer/libwasmer.so /usr/local/lib/
COPY ./suborbital /root/suborbital

ENV LD_LIBRARY_PATH=/usr/local/lib
ENV SUBO_DOCKER=1

WORKDIR /root/suborbital
RUN cargo build --target wasm32-wasi --lib --release

WORKDIR /root/runnable

ENTRYPOINT subo build --native .