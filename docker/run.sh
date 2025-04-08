#!/bin/bash
set -e

if [ ! -s "/var/lib/postgresql/data/PG_VERSION" ]; then
  echo "[+] Initializing PostgreSQL..."
  initdb -D /var/lib/postgresql/data
fi

echo "[+] Starting PostgreSQL in background..."
pg_ctl -D /var/lib/postgresql/data -w start

INIT_SQL="/opt/pandora/init.sql"
if [ -f "$INIT_SQL" ]; then
  echo "[+] Running database initialization script..."
  psql -U pandora -d postgres -f "$INIT_SQL"

  echo "[+] Launching pandora-core..."
  exec /usr/local/bin/pandora-core
else
  echo "[*] No init.sql script found at $INIT_SQL, skipping..."
fi
