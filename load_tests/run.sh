#!/bin/sh

docker run \
    -v $(pwd)/$1:/var/loadtest \
    -v $(pwd)/overload_token.txt:/var/loadtest/overload_token.txt \
    -v $SSH_AUTH_SOCK:/ssh-agent -e SSH_AUTH_SOCK=/ssh-agent \
    --net host \
    -it olgoncharov/ytank