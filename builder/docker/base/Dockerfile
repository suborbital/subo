ARG BUILDER_VERSION
FROM suborbital/subo:$BUILDER_VERSION as subo
FROM debian:bullseye-slim

WORKDIR /root/runnable

COPY --from=subo /go/bin/subo /usr/local/bin

RUN mkdir runnable; mkdir suborbital

ADD .env /root/.env
ADD install.sh /root/install.sh

