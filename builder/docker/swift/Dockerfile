FROM suborbital/subo:dev as subo

FROM ghcr.io/swiftwasm/swift:focal

COPY --from=subo /go/bin/subo /usr/local/bin

RUN mkdir -p /root/runnable

WORKDIR /root/runnable