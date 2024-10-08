#!/bin/bash

# ===========================================================
# your Temporal Cloud environment
# ===========================================================

export TEMPORAL_ACCOUNT="12345"

export TEMPORAL_TLS_CERT="$HOME/nexus-demo/certs/ca.pem"
export TEMPORAL_TLS_KEY="$HOME/nexus-demo/certs/ca.key"

export TEMPORAL_NAMESPACE_CALLER_BASE="my-caller-namespace"
export TEMPORAL_NAMESPACE_HANDLER_BASE="my-target-namespace"
export TEMPORAL_NAMESPACE_REGION="us-east-1"

export NEXUS_ENDPOINT="my-nexus-endpoint"

# ===========================================================
# likely don't need to change the stuff below
# ===========================================================

export TEMPORAL_ENV_SUBDOMAIN="tmprl"
export TEMPORAL_OPS_API="saas-api.${TEMPORAL_ENV_SUBDOMAIN}.cloud:443"

export TEMPORAL_NAMESPACE_CALLER="${TEMPORAL_NAMESPACE_CALLER_BASE}.${TEMPORAL_ACCOUNT}"
export TEMPORAL_NAMESPACE_HANDLER="${TEMPORAL_NAMESPACE_HANDLER_BASE}.${TEMPORAL_ACCOUNT}"

export TEMPORAL_ADDRESS_CALLER="${TEMPORAL_NAMESPACE_CALLER}.${TEMPORAL_ENV_SUBDOMAIN}.cloud:7233"
export TEMPORAL_ADDRESS_HANDLER="${TEMPORAL_NAMESPACE_HANDLER}.${TEMPORAL_ENV_SUBDOMAIN}.cloud:7233"
