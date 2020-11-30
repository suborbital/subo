FROM subo:dev as subo

FROM ghcr.io/swiftwasm/swift:focal

COPY --from=subo /go/bin/subo /usr/local/bin
COPY --from=subo /go/pkg/mod/github.com/wasmerio/wasmer-go@v0.3.2-0.20200903143934-299b1d478cde/wasmer/libwasmer.so /usr/local/lib/

ENV LD_LIBRARY_PATH=/usr/local/lib
ENV SUBO_DOCKER=1

RUN mkdir -p /root/runnable
WORKDIR /root/runnable

ENTRYPOINT subo build --native .