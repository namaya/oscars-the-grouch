#!/bin/bash

USERID=`http ':8080/api/users' name=namaya1 | jq -r .userId`

http ':8080/api/games' \
  Authorization:$USERID \
  name=game1

