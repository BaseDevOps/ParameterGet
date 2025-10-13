#!/bin/bash

BUILD_DIR=bin
NAME=aws-env

for GOOS in linux; do
    for GOARCH in arm64; do
        GOOS=$GOOS GOARCH=$GOARCH go build -v -o $BUILD_DIR/$NAME-$GOOS-$GOARCH
    done
done
