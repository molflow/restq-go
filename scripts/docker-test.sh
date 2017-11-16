#!/bin/bash

docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.8 /bin/bash -c 'go get github.com/jarcoal/httpmock && go test -v -cover'
