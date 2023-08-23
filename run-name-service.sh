#!/bin/sh -eu

REPOSITORY=registry.fillmore-labs.com
PROFILE=name-service

IMAGE=$(\
  env KO_DOCKER_REPO=$REPOSITORY/name-service \
  ko build --bare --sbom none ./cmd/name-service \
)

cat << _IMAGE > k8s/$PROFILE/image.yaml
---
apiVersion: builtin
kind: ImageTagTransformer
metadata:
  name: name-service
imageTag:
  name: name-service-image
  newName: ${IMAGE%%@*}
  digest: ${IMAGE#*@}
_IMAGE

kubectl apply -k k8s/$PROFILE
