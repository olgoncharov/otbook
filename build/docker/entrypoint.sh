#!/bin/sh

function migrate {
./bin/goose -dir ./migrations mysql "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?parseTime=true" up
}

NEXT_WAIT_TIME=1
until [ $NEXT_WAIT_TIME -eq 15 || migrate ]; do
    sleep $(( NEXT_WAIT_TIME++ ))
done

exec /bin/otbook