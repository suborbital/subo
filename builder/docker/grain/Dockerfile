FROM suborbital/subo:dev as subo

FROM ghcr.io/grain-lang/grain:0.4-slim
WORKDIR /root/module
COPY --from=subo /go/bin/subo /usr/local/bin/subo
RUN mkdir /root/suborbital
