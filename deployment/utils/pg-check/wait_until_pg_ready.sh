#!/bin/bash -e

until psql $STORAGE_POSTGRES_URI -c '\l'; do
  echo >&2 "$(date) Postgres is unavailable - sleeping"
  sleep 5
done
echo >&2 "$(date) Postgres is up"