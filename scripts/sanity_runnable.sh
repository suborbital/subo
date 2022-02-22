#!/bin/bash

set -eu

TEST_PROJECT="smoketest"

trap 'catch $? $LINENO' EXIT

catch() {
  if [[ "$1" != "0" ]]; then
    echo "Error $1 occurred on $2"
  fi

  # return to origin, clear directory stack
  pushd -0 > /dev/null && dirs -c

  if [[ -d "$TEST_PROJECT" ]]; then
    echo "cleaning up test artifacts"
    rm -rf "$TEST_PROJECT"
  fi
}

# create a new project
subo create project "$TEST_PROJECT"

# enter project directory
pushd "$TEST_PROJECT" > /dev/null

# create a runnable for each supported language
subo create runnable rs-test --lang rust
subo create runnable swift-test --lang swift
subo create runnable as-test --lang assemblyscript
subo create runnable tinygo-test --lang tinygo
subo create runnable js-test --lang javascript