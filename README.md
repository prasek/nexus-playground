# Temporal Nexus Playground

This is a playground for experimenting with [Temporal Nexus](https://temporal.io/nexus).

## Getting started

```
git clone https://github.com/prasek/nexus-playground
cd nexus-playground
```

## Get latest CLIs

```
brew install temporalio/brew/tcld
brew install temporal
```

Choose local dev or Temporal Cloud:
1. [Temporal Cloud Environment](#setup-temporal-cloud-environment)
1. [Local Dev Environment](#setup-temporal-local-dev-server-environment)


## Setup Temporal Cloud Environment

### Generate certs if you don't have them
```
mkdir -p $HOME/nexus-demo/certs
tcld gen ca --org temporal -d 1y --ca-cert $HOME/nexus-demo/certs/ca.pem --ca-key $HOME/nexus-demo/certs/ca.key
```

### Set envrionment config

edit [setEnv.sh](./setEnv.sh) and set the following to match your config:
```
export TEMPORAL_ACCOUNT="12345"

# Certs path for namespace creation and workers
export TEMPORAL_TLS_CERT="$HOME/nexus-demo/certs/ca.pem"
export TEMPORAL_TLS_KEY="$HOME/nexus-demo/certs/ca.key"

# Namespace without account suffix
export TEMPORAL_NAMESPACE_CALLER_BASE="my-caller-namespace"
export TEMPORAL_NAMESPACE_HANDLER_BASE="my-target-namespace"
export TEMPORAL_NAMESPACE_REGION="us-east-1"

# Nexus endpoint name that will be used
export NEXUS_ENDPOINT="my-nexus-endpoint"
```

### Create caller and handler namespaces
```
./cloud-namespaces.sh
```

### Create Nexus endpoint
```
./cloud-init.sh
```

### Run caller and handler workers

Handler worker:
```
./cloud-run.sh handler
```

Caller
```
./cloud-run.sh caller
```

### Browse to your Nexus Endpoint for instructions to get started
https://cloud.temporal.io/nexus

Done.

## Setup Temporal Local Dev Server Environment

Alternatively, create a local environment with the instructions below.

### Local dev server
```
./local-dev-server.sh
```

In the output, find the UI address, for example:
```
+ temporal server start-dev --dynamic-config-value system.enableNexus=true --http-port 7243
CLI 1.1.0 (Server 1.25.0, UI 2.30.3)

Server:  localhost:7233
UI:      http://localhost:8233
Metrics: http://localhost:60411/metrics
```

### Alternate: Use docker-compose
```
cd docker
docker-compose up
```

The UI address for this [docker-compose.yml](./docker/docker-compose.yml) is: http://localhost:8080

Note: when shutting down use the `--volumes` flag to get clean state:
```
docker-compose down --volumes
```

### Init
Create namespaces and a Nexus endpoint:
```
./local-init.sh
```

### Run caller and handler workers

Handler worker:
```
./local-run.sh handler
```

Caller worker:
```
./local-run.sh caller
```

### See additional instructions to get started
See [getting started](./service/description-local.md) and follow the instructions there.

Done.