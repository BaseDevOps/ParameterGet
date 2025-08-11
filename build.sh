#!/bin/bash

BUILD_DIR=bin
NAME=aws-env

for GOOS in linux; do
    for GOARCH in 386 amd64 arm64; do #
        echo "Building for $GOOS/$GOARCH..."
        GOOS=$GOOS GOARCH=$GOARCH go build -v -o $BUILD_DIR/$NAME-$GOOS-$GOARCH
    done
done

echo "Build complete."