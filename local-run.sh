#!/bin/bash

if [ $1 = "handler" ]; then
    ( set -x; cd handler; \
    go run ./worker \
        -target-host localhost:7233 \
        -namespace my-target-namespace
    )

elif [ $1 = "caller" ]; then
    ( set -x; cd caller; \
    go run ./worker \
        -target-host localhost:7233 \
        -namespace my-caller-namespace
    )

elif [ $1 = "starter" ]; then
    ( set -x; cd caller; \
    go run ./starter \
        -target-host localhost:7233 \
        -namespace my-caller-namespace \
        -endpoint my-nexus-endpoint \
        "${@:2}" \
    )

else
    echo "$1 not supported"
fi
