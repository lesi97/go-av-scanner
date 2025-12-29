#!/bin/sh
set -eu

mkdir -p /run/clamav /var/lib/clamav /var/log/clamav
chown -R clamav:clamav /run/clamav /var/lib/clamav /var/log/clamav

freshclam --config-file=/etc/clamav/freshclam.conf || true

clamd --config-file=/etc/clamav/clamd.conf &

freshclam --config-file=/etc/clamav/freshclam.conf -d &

exec /app/api
