FROM golang:1.18

RUN mkdir -p /root/github.com/suborbital/subo
WORKDIR /root/github.com/suborbital/subo

# dependencies first
COPY go.* ./
RUN go mod download

# then everything else
COPY subo ./subo
COPY builder ./builder
COPY deployer ./deployer
COPY packager ./packager
COPY publisher ./publisher
COPY project ./project
COPY scn ./scn

COPY *.go ./
COPY Makefile .

RUN make subo/docker-bin
