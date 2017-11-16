#!/bin/bash

docker run --rm -ti -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.8 bash -c "go get github.com/jarcoal/httpmock && bash"
