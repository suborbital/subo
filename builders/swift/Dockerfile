# build subo from scratch to support multi-arch builds
FROM golang:1.16 as subo

RUN mkdir -p /root/github.com/suborbital/subo
WORKDIR /root/github.com/suborbital/subo

COPY subo ./subo
COPY go.* .
COPY Makefile .

RUN make subo

FROM ghcr.io/swiftwasm/swift:focal

COPY --from=subo /go/bin/subo /usr/local/bin

ENV SUBO_DOCKER=1

RUN mkdir -p /root/runnable
WORKDIR /root/runnable

ENTRYPOINT subo build --native .