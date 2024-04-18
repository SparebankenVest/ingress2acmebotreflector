#!/bin/bash

# https://bertvv.github.io/cheat-sheets/Bash.html
set -o errexit   # abort on nonzero exitstatus
set -o nounset   # abort on unbound variable
set -o pipefail  # don't hide errors within pipes

CLUSTER_NAME="svai-charts-e2etest"

kind delete cluster -n "$CLUSTER_NAME"
