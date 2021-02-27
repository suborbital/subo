# build subo from scratch to support multi-arch builds
FROM golang:1.16 as subo

RUN mkdir -p /root/github.com/suborbital/subo
WORKDIR /root/github.com/suborbital/subo

COPY subo ./subo
COPY go.* .
COPY Makefile .

RUN make subo

FROM rust:1.50.0-slim-buster

WORKDIR /root
RUN mkdir runnable; mkdir suborbital

# install the wasm target and then install something that
# doesn't exist (and ignore the error) to update the crates.io index
RUN rustup target install wasm32-wasi
RUN cargo install lazy_static; exit 0

COPY --from=subo /go/bin/subo /usr/local/bin

ENV SUBO_DOCKER=1

WORKDIR /root/runnable

ENTRYPOINT subo build --native .
