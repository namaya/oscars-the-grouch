#!/bin/bash

http ':8080/api/users' name=bot | jq -r .userId
