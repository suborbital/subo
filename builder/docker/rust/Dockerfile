FROM suborbital/subo:dev as subo

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
