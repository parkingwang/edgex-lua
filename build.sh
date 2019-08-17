#!/usr/bin/env bash

GOOS=linux GOARCH=arm make image push
GOOS=linux GOARCH=arm64 make image push
GOOS=linux GOARCH=amd64 make image push
GOOS=linux make manifest