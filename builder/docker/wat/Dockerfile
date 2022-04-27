FROM suborbital/subo:dev as subo
FROM debian:bullseye as builder

WORKDIR /root

RUN apt-get update && \
    apt-get install pkg-config git build-essential libssl-dev clang cmake curl -y && \
    git clone -b 1.0.27 --recursive https://github.com/WebAssembly/wabt.git && \
    cd wabt && \
    git submodule update --init && \
    mkdir build && \
    cd build && \
    cmake .. && \
    cmake --build .

FROM debian:bullseye-slim

RUN mkdir runnable; mkdir suborbital

COPY --from=builder /root/wabt/bin/wat2wasm /usr/local/bin
COPY --from=subo /go/bin/subo /usr/local/bin

WORKDIR /root/runnable
