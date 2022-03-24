ARG BUILDER_VERSION
FROM suborbital/build-pack:${BUILDER_VERSION}

ENV CARGO_HOME=/usr/local/cargo\
    RUSTUP_HOME=/usr/local/rustup\
    PATH=/usr/local/cargo/bin:$PATH

RUN /root/install.sh rust