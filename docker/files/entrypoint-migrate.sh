#!/usr/bin/env sh

/wait-for-it.sh postgres-db-lb:5432 -t 600
cd /app
./migrate
