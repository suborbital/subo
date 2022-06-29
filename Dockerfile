FROM golang:1.18-bullseye AS builder
WORKDIR /root/github.com/suborbital/subo

COPY go.* ./
RUN go mod download

COPY cli ./cli
COPY builder ./builder
COPY deployer ./deployer
COPY packager ./packager
COPY publisher ./publisher
COPY project ./project
COPY scn ./scn
COPY *.go ./
COPY Makefile .
RUN make velo/docker-bin

FROM debian:bullseye
COPY --from=builder /go/bin/velo /usr/local/bin/velo
