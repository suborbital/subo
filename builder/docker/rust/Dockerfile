FROM suborbital/subo:dev as subo

FROM rust:1.56.1-slim-buster
WORKDIR /usr/src/runnable

# install the wasm target and then force an update of the crates.io index
RUN rustup target install wasm32-wasi && \
    cargo search suborbital

# make everything inside CARGO_HOME writable for anyone (again)
RUN chmod -R a+w $RUSTUP_HOME $CARGO_HOME

# embed subo
COPY --from=subo /go/bin/subo /usr/local/bin
