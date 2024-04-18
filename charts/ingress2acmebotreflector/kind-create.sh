#!/bin/bash

# https://bertvv.github.io/cheat-sheets/Bash.html
set -o errexit   # abort on nonzero exitstatus
set -o nounset   # abort on unbound variable
set -o pipefail  # don't hide errors within pipes

CLUSTER_NAME="svai-charts-e2etest"

echo "Creating cluster for end-to-end chart testing...🌟"

kind delete cluster -n "$CLUSTER_NAME"
kind create cluster --config cluster-config.yaml
kubectl config use-context "kind-${CLUSTER_NAME}"

echo "Installing infrastructure..."
helmfile init --force
helmfile apply

flux install --components-extra="image-reflector-controller,image-automation-controller"
echo "Installing infrastructure done"

echo "Cluster is ready to use 🎉"
