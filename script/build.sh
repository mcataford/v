#!/usr/bin/bash

OS=$(GOOS=$GOOS go env GOOS)
ARCH=$(GOARCH=$GOARCH go env GOARCH)

if [[ -z "$PROD" ]]; then
    CURRENT_COMMIT=$(git rev-parse HEAD)
    BUILD_PATH="build/$OS-$ARCH-$CURRENT_COMMIT"
    mkdir -p $BUILD_PATH
    GOOS=$OS GOARCH=$ARCH go build -o $BUILD_PATH -v -x .
else 
    BUILD_PATH="build/$OS-$ARCH"
    mkdir -p $BUILD_PATH .
    GOOS=$OS GOARCH=$ARCH go build -o $BUILD_PATH -ldflags "-w -s" .
fi
