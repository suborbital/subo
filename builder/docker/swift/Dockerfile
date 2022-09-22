FROM suborbital/subo:dev as subo

FROM ghcr.io/swiftwasm/swift:focal
WORKDIR /usr/src/runnable
COPY --from=subo /go/bin/subo /usr/local/bin
