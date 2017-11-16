#!/bin/bash

docker run -u $(id -u):$(id -g) -v /etc/passwd:/etc/passwd -v /etc/group:/etc/group --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.8 /bin/bash \
-c 'go get github.com/jarcoal/httpmock && \
GOPATH="$GOPATH:/usr/" go test -coverprofile cover.out && \
ln -s /usr/src/myapp /go/src/myapp && \
go tool cover -html=cover.out -o cover.html && \
rm /go/src/myapp' && \
sensible-browser cover.html &> /dev/null
