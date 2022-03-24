FROM ghcr.io/grain-lang/grain:0.4-slim

WORKDIR /root
RUN mkdir runnable; mkdir suborbital

COPY --from=suborbital/subo:dev /go/bin/subo /usr/local/bin

WORKDIR /root/runnable