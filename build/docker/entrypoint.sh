#!/bin/sh

function migrate {
./bin/goose -dir ./migrations mysql "${DB_USER_1}:${DB_PASSWORD_1}@tcp(${DB_HOST_1}:${DB_PORT_1})/${DB_NAME_1}?parseTime=true" up
}

NEXT_WAIT_TIME=1
until [[ $NEXT_WAIT_TIME -eq 15 ]] || migrate; do
    sleep $(( NEXT_WAIT_TIME++ ))
done

exec /bin/otbook