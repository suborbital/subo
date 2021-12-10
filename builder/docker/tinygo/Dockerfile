FROM suborbital/subo:dev as subo

FROM golang:bullseye as tinygobuilder
RUN apt-get update && apt-get install -y clang-11 llvm-11-dev lld-11 libclang-11-dev build-essential git cmake ninja-build

WORKDIR /root
RUN mkdir runnable
RUN git clone --depth 1 --branch v0.21.0 https://github.com/tinygo-org/tinygo.git
RUN git clone --depth 1 --branch version_103 https://github.com/WebAssembly/binaryen
RUN mkdir binaryen/build

WORKDIR /root/tinygo
RUN make llvm-source
RUN git submodule update --init --remote lib/wasi-libc
RUN export CC=clang-11 && export CXX=clang++-11 && make llvm-build
RUN make wasi-libc
RUN make

WORKDIR /root/binaryen/build
RUN cmake ../ -GNinja . && ninja wasm-opt
RUN cp -r bin lib /usr/local/.

COPY --from=subo /go/bin/subo /usr/local/bin

WORKDIR /root/runnable
