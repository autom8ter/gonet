#!/usr/bin/env bash

curl --header "Content-Type: application/json"   --request POST   --data '{"say":"hello"}'  localhost:8080/v1/echo