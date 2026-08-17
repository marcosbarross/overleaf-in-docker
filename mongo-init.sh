#!/usr/bin/env bash
set -euo pipefail

until mongosh --host mongo:27017 --quiet --eval "db.adminCommand('ping').ok" | grep 1 >/dev/null; do
  sleep 2
done

mongosh --host mongo:27017 --quiet <<'EOF'
try {
  const status = rs.status()
  if (status.ok === 1) {
    print("Replica set já inicializado.")
    quit(0)
  }
} catch (e) {}

rs.initiate({
  _id: "overleaf",
  members: [{ _id: 0, host: "mongo:27017" }]
})
EOF
