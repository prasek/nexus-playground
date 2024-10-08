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
echo "+ TEMPORAL_NAMESPACE_REGION=${TEMPORAL_NAMESPACE_REGION}"

set -x

tcld --server "${TEMPORAL_OPS_API}" login 

tcld --server "${TEMPORAL_OPS_API}" \
  namespace create \
	--namespace ${TEMPORAL_NAMESPACE_CALLER_BASE} \
	--region ${TEMPORAL_NAMESPACE_REGION} \
	--ca-certificate-file $TEMPORAL_TLS_CERT \
	--retention-days 30

tcld --server "${TEMPORAL_OPS_API}" \
  namespace create \
	--namespace ${TEMPORAL_NAMESPACE_HANDLER_BASE} \
	--region ${TEMPORAL_NAMESPACE_REGION} \
	--ca-certificate-file $TEMPORAL_TLS_CERT \
	--retention-days 30
