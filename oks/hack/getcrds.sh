#!/bin/sh

version=`octl kube kubectl foo get crd nodepools.oks.dev -o json | jq -r '.status.storedVersions[0]'`
mkdir -p api/v1beta2
octl kube kubectl foo get crd nodepools.oks.dev -o json | jq --arg version "$version" '{"openapi":"1.0.1","title":"nodepool","version":$version,"components":{"schemas":{"Nodepool": .spec.versions | map(select(.name == $version))[0].schema.openAPIV3Schema}}}' > api/oks.dev/docv1beta2/nodepool.yaml
