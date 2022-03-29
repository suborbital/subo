#!/bin/bash

# This script ensures that the version of TinyGo specified in the .tinygo-ver
# matches what is actually being built. This script should be run from the
# project root.

: ${GREP:=grep}

base=builder/docker/tinygo
version=$(cat ${base}/.tinygo-ver)

err_sum=0

function assert_file() {
    local file="$1"
    local expect="$2"
    local regexp="$3"

    local raw_match=$($GREP -noP "$regexp" ${file})
    local line_num=$(echo $raw_match | cut -f1 -d:)
    local match=$(echo $raw_match | cut -f2 -d:)

    if [ "$expect" == "$match" ]; then
        return 0
    else
        echo "ERROR at $file:$line_num: $expect != $match"
        err_sum=$(($err_sum + 1))
    fi
}

# Check that grep has Perl regexp support :^)
ver=$($GREP --version)
echo $ver | $GREP -q '(GNU grep'
if [ $? -ne 0 ]; then
    if ! command -v ggrep &> /dev/null; then
        echo "ERROR: Must use GNU grep"
        if [ $(uname) == "Darwin" ]; then
            echo "  brew install grep"
        fi
        exit 1
    else
        GREP=ggrep
    fi
fi

# Check the base image
assert_file $base/Dockerfile.base $version 'branch \K.+(?= https://github.com/tinygo-org/tinygo.git)'
assert_file $base/Dockerfile $version 'FROM suborbital/tinygo-base:\K.+(?= as)'

if [ $err_sum -eq 0 ]; then
    echo "Success"
else
    echo "Failed with $err_sum errors"
fi

exit $err_sum
