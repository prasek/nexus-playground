#!/bin/bash

set -x

temporal operator namespace create --namespace my-target-namespace
temporal operator namespace create --namespace my-caller-namespace

temporal operator nexus endpoint delete --name myendpoint
temporal operator nexus endpoint create \
  --name my-nexus-endpoint \
  --target-namespace my-target-namespace \
  --target-task-queue my-handler-task-queue \
  --description-file ./service/description.md
