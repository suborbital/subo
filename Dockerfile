FROM golang:1.18-bullseye AS builder
WORKDIR /root/github.com/suborbital/subo

COPY go.* ./
RUN go mod download

COPY subo ./subo
COPY builder ./builder
COPY deployer ./deployer
COPY packager ./packager
COPY publisher ./publisher
COPY project ./project
COPY se2 ./se2
COPY *.go ./
COPY Makefile .
RUN make subo/docker-bin

FROM debian:bullseye
COPY --from=builder /go/bin/subo /go/bin/subo
