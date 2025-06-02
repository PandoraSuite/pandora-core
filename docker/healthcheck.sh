#!/bin/sh

if ! wget --spider --quiet http://localhost:80/health; then
  echo "REST healthcheck failed"
  exit 1
fi

if command -v grpc_health_probe > /dev/null 2>&1; then
  if ! grpc_health_probe -addr=localhost:50051; then
    echo "gRPC healthcheck failed"
    exit 1
  fi
else
  echo "grpc_health_probe not found, skipping gRPC check"
fi

exit 0
