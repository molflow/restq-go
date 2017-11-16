#!/bin/bash

./restq -c > epa
cat epa | ./restq -f epa -p
if [ "$1" == "-v" ]; then
    cat epa
    ./restq -f epa -g | jq -r .queue
    cat epa | ./restq -f epa -p
fi
if [ $(./restq -f epa -g | jq -r .queue) == $(jq -r .queue epa) ]; then
    rm epa
    echo OK
    exit 0
else
    rm epa
    echo FAIL
    exit 1
fi

