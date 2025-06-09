#!/bin/bash
TMPFILE=$(mktemp)

# first create a tmpfile to store the results from the psql command 
connected=false
trap 'rm -f "$TMPFILE"' EXIT
#ensure cleanup after operation


for i in {0..60}; do 
PGPASSWORD="postgres" psql -h localhost -U postgres -p 5432 -c '\q' 2> /dev/null
if [ $? -eq 0 ]; then
    connected=true
    echo "docker container has initialised"
    break
else
    echo "waiting for postgres to start: ($i)"
    sleep 0.5
fi
done

if [ "$connected" = false ]; then
  echo "Postgres did not start in time"
  exit 1
fi
# attempt to connect to the database

PGPASSWORD="postgres" psql -h localhost -U postgres -p 5432 -c "CREATE DATABASE gator" 2> "$TMPFILE"
#attempt to create the gator database

if grep -q "already exists" "$TMPFILE"; then
    echo "Database already exists - Continuing ..."
#error will be specific in the event the gator database already exists

elif [ -s "$TMPFILE" ];
then
    echo "Some other error occured."
    cat "$TMPFILE"
    echo "error occured in database creation: skipping to closing container"
    exit 1
# -s will test to see if tmpfile exists, and whether the contents are larger than zero
else 
    echo "gator database created successfully"
fi

echo "Running migrations"

