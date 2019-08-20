#!/usr/bin/env bash

SUDO="sudo "
if [ "Darwin" == "$(uname -s)" ]; then
    SUDO=""
fi

OS_SUDO=${SUDO} GOOS=linux GOARCH=arm OSARCH=arm32v7 make image push
OS_SUDO=${SUDO} GOOS=linux GOARCH=arm64 OSARCH=arm64v8 make image push
OS_SUDO=${SUDO} GOOS=linux GOARCH=amd64 OSARCH=amd64 make image push
OS_SUDO=${SUDO} GOOS=linux make manifest