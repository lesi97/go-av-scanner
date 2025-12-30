#!/bin/sh
set -eu

mkdir -p /run/clamav /var/lib/clamav /var/log/clamav
chown -R clamav:clamav /run/clamav /var/lib/clamav /var/log/clamav

freshclam --config-file=/etc/clamav/freshclam.conf || true

clamd --config-file=/etc/clamav/clamd.conf &

for i in $(seq 1 30); do
  if clamdscan --no-summary --fdpass /etc/hosts >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

freshclam --config-file=/etc/clamav/freshclam.conf -d &

exec /app/api
