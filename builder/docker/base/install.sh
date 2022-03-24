#!/bin/bash

# Exit on error or if there are any unbounded variables
set -eu

trap 'catch $? $LINENO' EXIT

catch() {
  if [[ "$1" != "0" ]]; then
    echo "An Error $1 occurred on $2"
  fi

  # return to origin, clear directory stack
  pushd -0 > /dev/null && dirs -c
}

# push source root directory to front of directory stack
pushd "$( pushd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
source .env

function prepare_rust {
  apt-get install build-essential -y

  dpkgArch="$(dpkg --print-architecture)"
  case "${dpkgArch##*-}" in
    amd64) rustArch='x86_64-unknown-linux-gnu' ;;
    arm64) rustArch='aarch64-unknown-linux-gnu' ;;
    *) echo >&2 "unsupported architecture: ${dpkgArch}"; exit 1 ;;
  esac;

  wget "https://static.rust-lang.org/rustup/archive/1.24.3/${rustArch}/rustup-init"
  chmod +x rustup-init
  ./rustup-init -y --no-modify-path --profile minimal --default-toolchain ${RUST_VERSION} --default-host ${rustArch}
  rm rustup-init

  rustup target install wasm32-wasi

  cargo install lazy_static || echo "cargo index preloaded...probably"
}

function prepare_javy {
  # javy extends rust image which has already cleaned up it's indexes
  apt-get update
  apt-get install cmake git clang-11 nodejs npm -y
  git clone -b suborbital-v0.2.0 https://github.com/suborbital/javy.git --single-branch
  make -C javy
}

function prepare_javascript {
  apt-get install xz-utils -y
  dpkgArch="$(dpkg --print-architecture)"
  case "${dpkgArch##*-}" in
    amd64) nodeArch='x64';;
    arm64) nodeArch='arm64';;
    *) echo >&2 "unsupported architecture: ${dpkgArch}"; exit 1 ;;
  esac
  curl -fsSLO --compressed "https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-${nodeArch}.tar.xz"
  tar -xJf "node-v${NODE_VERSION}-linux-${nodeArch}.tar.xz" -C /usr/local --strip-components=1 --no-same-owner
  rm "node-v${NODE_VERSION}-linux-${nodeArch}.tar.xz"
  ln -s /usr/local/bin/node /usr/local/bin/nodejs
  chmod -R o=u /root
  apt purge xz-utils -y
}

function prepare_tinygo {
  dpkgArch="$(dpkg --print-architecture)"

  # Install golang
  wget "https://go.dev/dl/go${GOLANG_VERSION}.linux-${dpkgArch}.tar.gz"
  tar -C /usr/local -xzf "go${GOLANG_VERSION}.linux-${dpkgArch}.tar.gz"
  rm "go${GOLANG_VERSION}.linux-${dpkgArch}.tar.gz"

  # Install tinygo
  wget "https://github.com/tinygo-org/tinygo/releases/download/v${TINYGO_VERSION}/tinygo_${TINYGO_VERSION}_${dpkgArch}.deb" &&\
  dpkg -i tinygo_${TINYGO_VERSION}_${dpkgArch}.deb
  rm tinygo_${TINYGO_VERSION}_${dpkgArch}.deb

  ln -s /usr/local/go/bin /usr/local/bin/go
  go mod download github.com/suborbital/reactr@latest
}


if [ "$#" -ne 1 ]; then
  echo "Usage: $0 profile" >&2
  exit 1
fi

apt update && apt upgrade -y
apt-get install ca-certificates wget curl --no-install-recommends -y

case "${1##*-}" in
  rust) prepare_rust ;;
  javy) prepare_javy ;;
  javascript) prepare_javascript ;;
  tinygo) prepare_tinygo ;;
  *) echo "Unsupported language" >&2; exit 1 ;;
esac

apt purge wget curl -y
#apt-mark auto '.*' > /dev/null
apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false
rm -rf /var/lib/apt/lists/*
# fix up anything that may have broken during cleanup
apt -f install
rm -rf /usr/local/share/man
rm -rf /usr/local/share/doc