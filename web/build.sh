#!/usr/bin/env bash
# Build the web image directly into minikube's Docker daemon.
set -euo pipefail
cd "$(dirname "$0")"
IMAGE="${IMAGE:-kubby-manatee-web:latest}"
eval "$(minikube docker-env)"
DOCKER_BUILDKIT=0 docker build -t "$IMAGE" .
