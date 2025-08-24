#!/bin/sh

echo "[HEALTHCHECK] Starting health verification..."

# Check REST API
echo "[HEALTHCHECK] Checking REST API (port 80)..."
if ! curl -f -s http://localhost:80/api/v1/health; then
  echo "[HEALTHCHECK] FAILED: REST API is not responding"
  exit 1
fi
echo "[HEALTHCHECK] REST API: OK"

# Check gRPC
echo "[HEALTHCHECK] Checking gRPC service (port 50051)..."
if command -v grpc_health_probe > /dev/null 2>&1; then
  if ! grpc_health_probe -addr=localhost:50051 > /dev/null 2>&1; then
    echo "[HEALTHCHECK] FAILED: gRPC service is not responding"
    exit 1
  fi
  echo "[HEALTHCHECK] gRPC service: OK"
else
  echo "[HEALTHCHECK] WARNING: grpc_health_probe not found, skipping gRPC check"
fi

echo "[HEALTHCHECK] All services healthy"
exit 0
