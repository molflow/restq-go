#!/bin/bash

./gogets.sh
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build --ldflags '-extldflags "-static"'
