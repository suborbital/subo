ARG BUILDER_VERSION
FROM suborbital/javy:v${BUILDER_VERSION} as javy
FROM suborbital/build-pack:${BUILDER_VERSION}

COPY --from=javy /root/javy/target/release/javy /usr/local/bin

RUN /root/install.sh javascript

#RUN chmod -R o=u /root