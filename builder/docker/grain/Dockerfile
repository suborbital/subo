FROM suborbital/subo:dev as subo

FROM ghcr.io/grain-lang/grain:0.5-slim
WORKDIR /root/runnable
COPY --from=subo /go/bin/subo /usr/local/bin/subo
# Necessary to build Grain runnables
COPY --from=subo /usr/bin/curl /usr/bin/curl
RUN mkdir /root/suborbital
