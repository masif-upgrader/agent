#!/bin/bash
set -e

/usr/sbin/addgroup --system masif-upgrader-agent
exec /bin/systemctl daemon-reload
