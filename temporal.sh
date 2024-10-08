#!/bin/bash

if [ -f ./setEnv.sh ]; then
    echo sourced from setEnv.sh, delete or modify to use alternate:
    source ./setEnv.sh
fi

if [ ${ENV} = "handler" ]; then
    export TEMPORAL_ADDRESS=${TEMPORAL_ADDRESS_HANDLER}
    export TEMPORAL_NAMESPACE=${TEMPORAL_NAMESPACE_HANDLER}
else
    export TEMPORAL_ADDRESS=${TEMPORAL_ADDRESS_CALLER}
    export TEMPORAL_NAMESPACE=${TEMPORAL_NAMESPACE_CALLER}
fi

echo "+ TEMPORAL_ADDRESS=${TEMPORAL_ADDRESS}"
echo "+ TEMPORAL_NAMESPACE=${TEMPORAL_NAMESPACE}"
echo "+ TEMPORAL_TLS_CERT=${TEMPORAL_TLS_CERT}"
echo "+ TEMPORAL_TLS_KEY=${TEMPORAL_TLS_KEY}"

set -x

temporal "${@:1}"