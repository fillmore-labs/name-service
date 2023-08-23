#!/bin/sh -eu

env KO_DOCKER_REPO=registry.fillmore-labs.com/name-client \
  ko run ./cmd/name-client \
  --bare --sbom none \
  -- -n names --env="NAME_SERVICE=name-service.names.svc.cluster.local:80"
