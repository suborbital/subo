FROM subo:dev as subo

FROM ghcr.io/swiftwasm/swift:focal

COPY --from=subo /go/bin/subo /usr/local/bin

ENV SUBO_DOCKER=1

RUN mkdir -p /root/runnable
WORKDIR /root/runnable

ENTRYPOINT subo build --native .