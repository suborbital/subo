FROM localhost:5000/suborbital/subo:dev as subo

FROM localhost:5000/grain-lang/grain/grain:main

WORKDIR /root
RUN mkdir runnable; mkdir suborbital

COPY --from=subo /go/bin/subo /usr/local/bin

ENV SUBO_DOCKER=1

WORKDIR /root/runnable

ENTRYPOINT subo build --native .
