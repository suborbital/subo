FROM golang:1.15 as builder

RUN mkdir -p /root/github.com/suborbital/subo
WORKDIR /root/github.com/suborbital/subo

COPY subo ./subo
COPY go.* .
COPY Makefile .

RUN make subo