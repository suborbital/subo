FROM ghcr.io/swiftwasm/swift:focal

COPY --from=suborbital/subo:dev  /go/bin/subo /usr/local/bin

RUN mkdir -p /root/runnable

WORKDIR /root/runnable