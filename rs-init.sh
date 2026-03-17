#!/bin/bash
# Initiates the replica set if not yet done, then blocks until PRIMARY.

set -e

MONGO_HOST="localhost:27017"
MONGO_USER="admin"
MONGO_PASS="password"

echo "Waiting for mongod to accept connections..."
until mongosh --host "$MONGO_HOST" \
              -u "$MONGO_USER" -p "$MONGO_PASS" \
              --authenticationDatabase admin \
              --quiet --eval "db.adminCommand('ping')" > /dev/null 2>&1; do
  sleep 1
done

echo "mongod is up. Checking replica set status..."
STATUS=$(mongosh --host "$MONGO_HOST" \
                 -u "$MONGO_USER" -p "$MONGO_PASS" \
                 --authenticationDatabase admin \
                 --quiet --eval "rs.status().ok" 2>/dev/null || echo "0")

if [ "$STATUS" != "1" ]; then
  echo "Initiating replica set rs0..."
  mongosh --host "$MONGO_HOST" \
          -u "$MONGO_USER" -p "$MONGO_PASS" \
          --authenticationDatabase admin \
          --quiet --eval '
    rs.initiate({
      _id: "rs0",
      members: [{ _id: 0, host: "mongo:27017" }]
    })
  '
  echo "Replica set initiated."
else
  echo "Replica set already initiated."
fi

echo "Waiting for PRIMARY election..."
until mongosh --host "$MONGO_HOST" \
              -u "$MONGO_USER" -p "$MONGO_PASS" \
              --authenticationDatabase admin \
              --quiet --eval "db.isMaster().ismaster" 2>/dev/null | grep -q "true"; do
  sleep 1
done

echo "Node is PRIMARY. Replica set ready."