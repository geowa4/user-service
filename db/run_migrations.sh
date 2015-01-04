#!/bin/bash

for f in /docker-entrypoint-initdb.d/migrations/*.sql; do
  [ -f "$f" ] && cat $f | gosu postgres postgres --single user-service -j
done
