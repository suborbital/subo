FROM golang:1.17

RUN mkdir -p /root/github.com/suborbital/subo
WORKDIR /root/github.com/suborbital/subo

# dependencies first
COPY go.* ./
RUN go mod download

# then everything else
COPY subo ./subo
COPY builder ./builder
COPY scn ./scn
COPY *.go ./
COPY Makefile .

RUN make subo