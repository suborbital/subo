FROM suborbital/subo:dev as subo

FROM ghcr.io/swiftwasm/swift:focal
WORKDIR /root/module
COPY --from=subo /go/bin/subo /usr/local/bin
