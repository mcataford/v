#!/usr/bin/bash

if [[ -z "$FIX" ]]; then
    [ -z "$(gofmt -s -l ./)" ] || exit 1
else
    gofmt -s -l -w ./
fi
