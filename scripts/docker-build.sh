#!/bin/bash

docker run -u $(id -u):$(id -g) -v /etc/passwd:/etc/passwd -v /etc/group:/etc/group --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.8 ./scripts/build_linux.sh
