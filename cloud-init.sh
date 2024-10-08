#!/bin/bash

if [ -f ./setEnv.sh ]; then
    echo sourced from setEnv.sh, delete or modify to use alternate:
    source ./setEnv.sh
fi

echo "+ NEXUS_ENDPOINT=${NEXUS_ENDPOINT}"
echo "+ TEMPORAL_ADDRESS_CALLER=${TEMPORAL_ADDRESS_CALLER}"
echo "+ TEMPORAL_NAMESPACE_CALLER=${TEMPORAL_NAMESPACE_CALLER}"
echo "+ TEMPORAL_ADDRESS_HANDLER=${TEMPORAL_ADDRESS_HANDLER}"
echo "+ TEMPORAL_NAMESPACE_HANDLER=${TEMPORAL_NAMESPACE_HANDLER}"
echo "+ TEMPORAL_TLS_CERT=${TEMPORAL_TLS_CERT}"
echo "+ TEMPORAL_TLS_KEY=${TEMPORAL_TLS_KEY}"
echo "+ TEMPORAL_OPS_API=${TEMPORAL_OPS_API}"

set -x

tcld --server "${TEMPORAL_OPS_API}" \
  nexus endpoint delete \
  --name ${NEXUS_ENDPOINT} \

until tcld --server "${TEMPORAL_OPS_API}" \
  nexus endpoint create \
  --name ${NEXUS_ENDPOINT} \
  --target-task-queue my-handler-task-queue \
  --target-namespace ${TEMPORAL_NAMESPACE_HANDLER} \
  --allow-namespace ${TEMPORAL_NAMESPACE_CALLER} \
  --description-file ./service/description.md
do
    sleep 1
    echo "Trying again ..."
done
