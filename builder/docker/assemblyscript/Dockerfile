FROM suborbital/subo:dev as subo

FROM node:16-buster-slim

RUN npm install -g npm@latest

COPY --from=subo /go/bin/subo /usr/local/bin

ENV SUBO_DOCKER=1

WORKDIR /root/runnable

ENTRYPOINT subo build --native --langs assemblyscript .