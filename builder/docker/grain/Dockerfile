FROM suborbital/subo:dev as subo

FROM ghcr.io/grain-lang/grain:0.5-slim
WORKDIR /root/runnable
COPY --from=subo /go/bin/subo /usr/local/bin/subo

# Necessary to build Grain runnables
RUN apt-get update && apt-get install curl -y && rm -rf /var/lib/apt/lists/*

RUN mkdir /root/suborbital
