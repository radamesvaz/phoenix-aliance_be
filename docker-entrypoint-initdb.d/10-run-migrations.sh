#!/bin/sh
set -eu

# This script is executed by the official Postgres Docker image ONLY on first init
# (i.e., when the data directory is empty).
#
# We mount migrations to /migrations and apply ONLY *_up.sql, so *_down.sql is never executed here.

echo "Running SQL migrations (*_up.sql) from /migrations ..."

if [ ! -d "/migrations" ]; then
  echo "ERROR: /migrations directory not found inside container."
  exit 1
fi

# Apply in name order (001, 002, ...)
# shellcheck disable=SC2039
for file in /migrations/*_up.sql; do
  # If no files match, the glob remains literal in POSIX sh; guard it.
  if [ ! -f "$file" ]; then
    break
  fi
  echo "Applying: $(basename "$file")"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f "$file"
done

echo "Migrations completed."
