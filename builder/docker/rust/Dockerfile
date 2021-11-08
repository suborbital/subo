FROM suborbital/subo:dev as subo

FROM rust:1.56.1-slim-buster

WORKDIR /root
RUN mkdir runnable; mkdir suborbital

# install the wasm target and then install something that
# doesn't exist (and ignore the error) to update the crates.io index
RUN rustup target install wasm32-wasi
RUN cargo install lazy_static; exit 0

COPY --from=subo /go/bin/subo /usr/local/bin

WORKDIR /root/runnable