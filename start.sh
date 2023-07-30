#!/bin/sh

set -e

echo "run db migrations"
/app/soda migrate

echo "start the app"
exec "$@"