FROM suborbital/tinygo-base:v0.22.0 as builder
FROM suborbital/subo:dev as subo

FROM golang:1.17-bullseye
WORKDIR /root/runnable
COPY --from=subo /go/bin/subo /usr/local/bin
COPY --from=builder /root/tinygo/build/release.tar.gz /usr/local/.
RUN cd /usr/local && \
    tar -xzf release.tar.gz && rm release.tar.gz && \
    ln -s /usr/local/tinygo/bin/tinygo /usr/local/bin/tinygo && \
    ln -s /usr/local/tinygo/bin/wasm-opt /usr/local/bin/wasm-opt
RUN go mod download github.com/suborbital/reactr@latest
