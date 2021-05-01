# build subo from scratch to support multi-arch builds
FROM golang:1.16 as subo

RUN mkdir -p /root/github.com/suborbital/subo
WORKDIR /root/github.com/suborbital/subo

COPY subo ./subo
COPY go.* .
COPY Makefile .

RUN make subo

FROM node:16-buster-slim

RUN npm install -g npm@latest

COPY --from=subo /go/bin/subo /usr/local/bin

ENV SUBO_DOCKER=1

WORKDIR /root/runnable

ENTRYPOINT subo build --native .