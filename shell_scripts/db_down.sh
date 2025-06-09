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
# attempt to connect to the database


echo "running down migration"
cd "./sql/schema"
CONNSTRING="postgres://postgres:postgres@localhost:5432/gator"
goose postgres "$CONNSTRING" down 2>"$TMPFILE"
# attempt to run the migration

elif [ -s "$TMPFILE" ]; then
    echo "error occured when running down migration"
    cat "$TMPFILE"
else 
    echo "down migration ran successfully"
fi

