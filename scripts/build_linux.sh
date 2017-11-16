#!/bin/bash

CGO_ENABLED=0 go build --ldflags '-extldflags "-static"' restq.go
