#!/bin/bash

#!/bin/bash

if [ -f ./setEnv.sh ]; then
    echo sourced from setEnv.sh, delete or modify to use alternate:
    source ./setEnv.sh
fi

#echo "+ NEXUS_ENDPOINT=${NEXUS_ENDPOINT}"
#echo "+ TEMPORAL_ADDRESS_CALLER=${TEMPORAL_ADDRESS_CALLER}"
#echo "+ TEMPORAL_NAMESPACE_CALLER=${TEMPORAL_NAMESPACE_CALLER}"
#echo "+ TEMPORAL_ADDRESS_HANDLER=${TEMPORAL_ADDRESS_HANDLER}"
#echo "+ TEMPORAL_NAMESPACE_HANDLER=${TEMPORAL_NAMESPACE_HANDLER}"
#echo "+ TEMPORAL_TLS_CERT=${TEMPORAL_TLS_CERT}"
#echo "+ TEMPORAL_TLS_KEY=${TEMPORAL_TLS_KEY}"

if [ $1 = "handler" ]; then
    ( set -x; cd handler; \
    go run ./worker \
        -target-host $TEMPORAL_ADDRESS_HANDLER \
        -namespace $TEMPORAL_NAMESPACE_HANDLER \
        -client-cert $TEMPORAL_TLS_CERT \
        -client-key $TEMPORAL_TLS_KEY  \
    )

elif [ $1 = "caller" ]; then
    ( set -x; cd caller; \
    go run ./worker \
        -target-host $TEMPORAL_ADDRESS_CALLER \
        -namespace $TEMPORAL_NAMESPACE_CALLER \
        -client-cert $TEMPORAL_TLS_CERT \
        -client-key $TEMPORAL_TLS_KEY  \
    )

elif [ $1 = "starter" ]; then
    ( set -x; cd caller; \
    go run ./starter \
        -target-host $TEMPORAL_ADDRESS_CALLER \
        -namespace $TEMPORAL_NAMESPACE_CALLER \
        -client-cert $TEMPORAL_TLS_CERT \
        -client-key $TEMPORAL_TLS_KEY  \
        -endpoint $NEXUS_ENDPOINT \
        "${@:2}" \
    )

else
    echo "$1 not supported"
fi
