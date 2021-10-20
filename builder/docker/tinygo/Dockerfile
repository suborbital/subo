FROM suborbital/subo:dev as subo

# FROM tinygo/tinygo:0.20.0
# doesn't work on M1 :(

FROM golang:bullseye as tinygobuilder
RUN apt update && apt install -y clang-11 llvm-11-dev lld-11 libclang-11-dev build-essential git cmake ninja-build

WORKDIR /root
RUN mkdir runnable; mkdir suborbital
RUN git clone https://github.com/tinygo-org/tinygo.git
WORKDIR /root/tinygo
RUN git submodule update --init --remote lib/wasi-libc
RUN make wasi-libc
RUN CGO_ENABLED=1 go install

COPY --from=subo /go/bin/subo /usr/local/bin

ENV SUBO_DOCKER=1

WORKDIR /root

# temporary hack
RUN git clone https://github.com/suborbital/reactr
WORKDIR /root/reactr
RUN git checkout jagger/tinygo

WORKDIR /root/runnable

ENTRYPOINT subo build --native .
