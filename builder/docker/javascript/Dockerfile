FROM suborbital/subo:dev as subo

FROM node:16-bullseye-slim

WORKDIR /root
RUN mkdir runnable; mkdir suborbital

ENV PATH="/root/.cargo/bin:$PATH"

RUN apt-get update && \
    apt-get install pkg-config git build-essential libssl-dev clang cmake curl -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y && \
    rustup target install wasm32-wasi && \
    cargo install wasmtime-cli cargo-wasi && \
    git clone https://github.com/suborbital/javy.git && \
    cd javy && \
    make && cargo install --path crates/cli && \
    cd - && \
    rm -rf javy && \
    mv /root/.cargo/bin/javy /usr/local/bin && \
    rustup self uninstall -y && \
    apt-get --purge remove pkg-config git build-essential libssl-dev clang cmake curl -y

COPY --from=subo /go/bin/subo /usr/local/bin

WORKDIR /root/runnable
