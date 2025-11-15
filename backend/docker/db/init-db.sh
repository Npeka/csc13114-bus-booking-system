#!/bin/bash
set -e

DB_LIST="bus_dev_user bus_dev_payment bus_dev_booking"

for db in $DB_LIST; do
  echo "Creating database '$db'"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE DATABASE "$db";
EOSQL
done
