#!/usr/bin/env bash

OSARCH=arm make $*
OSARCH=amd64 make $*
