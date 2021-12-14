FROM suborbital/tinygo-base:v0.21.0 as builder
FROM suborbital/subo:dev as subo

FROM golang:bullseye

COPY --from=subo /go/bin/subo /usr/local/bin
COPY --from=builder /root/tinygo/build/release.tar.gz /usr/local/.
WORKDIR /usr/local
RUN tar -xzf release.tar.gz \
    && rm release.tar.gz

RUN ln -s /usr/local/tinygo/bin/tinygo /usr/local/bin/tinygo
RUN ln -s /usr/local/tinygo/bin/wasm-opt /usr/local/bin/wasm-opt

WORKDIR /root/runnable

RUN go mod download github.com/suborbital/reactr@latest
