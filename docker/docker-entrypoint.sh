#!/bin/sh
set -e

export POSTGRES_USER="$PANDORA_DB_USER"
export POSTGRES_PASSWORD="$PANDORA_DB_PASSWORD"
export POSTGRES_DB="$PANDORA_DB_NAME"

export PANDORA_DB_DNS="host=localhost port=5432 user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB sslmode=disable timezone=UTC"

echo "[INFO] Launching PostgreSQL..."

/usr/local/bin/docker-entrypoint.sh postgres > /var/log/pandora/postgres.log 2>&1 &

PG_PID=$!

if ! kill -0 "$PG_PID" 2>/dev/null; then
  echo "[ERROR] PostgreSQL failed to start."
  exit 1
fi

echo "[INFO] Waiting for PostgreSQL to be ready..."

until pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB" -h localhost > /dev/null 2>&1; do
  sleep 1
done

echo "[INFO] PostgreSQL is ready."

echo "[INFO] Starting Pandora Core..."

exec /usr/local/bin/pandora-core
