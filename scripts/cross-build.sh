#!/bin/bash
BINARY_NAME=goxz

for GOOS in linux darwin
do
	for GOARCH in amd64
		do

		echo Building $GOOS/$GOARCH
		GOOS=$GOOS GOARCH=$GOARCH go build -o build/$BINARY_NAME-$GOOS-$GOARCH .
    zip build/$BINARY_NAME-$GOOS-$GOARCH.zip build/$BINARY_NAME-$GOOS-$GOARCH
	done
done
