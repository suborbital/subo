#!/bin/bash

set -eu

TEST_PROJECT="buildertest"

trap 'catch $? $LINENO' EXIT

catch() {
  if [[ "$1" != "0" ]]; then
    echo "An Error $1 occurred on $2"
  fi

  # return to origin, clear directory stack
  pushd -0 > /dev/null && dirs -c

  if [[ -d "$TEST_PROJECT" ]]; then
    echo "Cleaning up test artifacts..."
    rm -rf "$TEST_PROJECT" || echo "Failed to clean up test artifacts, if this was a permissions error try using 'sudo rm -rf $TEST_PROJECT'"
  fi
}

# build local subo docker image
docker image ls -q suborbital/subo:dev
if [[ "$?" != "0" ]]; then
  make subo/docker
fi

# rebuild local docker build tooling
docker build . -f builder/docker/assemblyscript/Dockerfile -t suborbital/builder-as:dev
docker build . -f builder/docker/grain/Dockerfile --platform linux/amd64 -t suborbital/builder-gr:dev
docker build . -f builder/docker/javascript/Dockerfile -t suborbital/builder-js:dev
docker build . -f builder/docker/rust/Dockerfile -t suborbital/builder-rs:dev
docker build . -f builder/docker/swift/Dockerfile -t suborbital/builder-swift:dev
docker build . -f builder/docker/tinygo/Dockerfile --platform linux/amd64 --build-arg TARGETARCH=amd64 -t suborbital/builder-tinygo:dev
docker build . -f builder/docker/wat/Dockerfile -t suborbital/builder-wat:dev


# create a new project
subo create project "$TEST_PROJECT"

# enter project directory
pushd "$TEST_PROJECT" > /dev/null

# create a runnable for each supported language
subo create module rs-test --lang rust
subo create module swift-test --lang swift
subo create module as-test --lang assemblyscript
subo create module tinygo-test --lang tinygo
subo create module js-test --lang javascript

# build project bundle
subo build . --builder-tag dev
