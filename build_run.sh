#!/usr/bin/env bash

docker build -t bentsolheim/kilsundvaeret-api .
docker run \
 --rm \
 -p 9010:9010 \
 --name kilsundvaeret-api \
 bentsolheim/kilsundvaeret-api:latest
