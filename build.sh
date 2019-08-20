#!/usr/bin/env bash

SUDO="sudo "
if [ "Darwin" == "$(uname -s)" ]; then
    SUDO=""
fi

OS_SUDO=${SUDO} GOOS=linux GOARCH=arm make image push
OS_SUDO=${SUDO} GOOS=linux GOARCH=arm64 make image push
OS_SUDO=${SUDO} GOOS=linux GOARCH=amd64 make image push
OS_SUDO=${SUDO} GOOS=linux make manifest