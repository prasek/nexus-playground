#!/bin/bash

set -x

#curl -sSf https://temporal.download/cli.sh | sh -s -- --version v1.1.0 --dir .

temporal server start-dev --dynamic-config-value system.enableNexus=true --http-port 7243
