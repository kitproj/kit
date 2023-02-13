#!/usr/bin/env bash
set -eu -o pipefail

kubectl -n argo port-forward "svc/minio" "9000:9000" &

