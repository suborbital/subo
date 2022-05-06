FROM suborbital/tinygo-base:v0.23.0 as tinygo
FROM suborbital/subo:dev as subo

FROM golang:1.18-bullseye
WORKDIR /root/runnable
COPY --from=tinygo /root/release/tinygo /usr/local/tinygo
COPY --from=subo /go/bin/subo /usr/local/bin
RUN ln -s /usr/local/tinygo/bin/tinygo /usr/local/bin/tinygo && \
    ln -s /usr/local/tinygo/bin/wasm-opt /usr/local/bin/wasm-opt
RUN go mod download github.com/suborbital/reactr@latest
