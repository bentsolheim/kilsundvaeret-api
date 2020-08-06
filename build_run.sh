#!/usr/bin/env bash

docker build -t bentsolheim/kilsundvaeret-api .
docker run \
 --rm \
 -p 8080:9010 \
 --name kilsundvaeret-api \
 --env DB_HOST=kilsundvaeret-api_db_1 \
 --env DATA_RECEIVER_URL=http://something.no \
 --net kilsundvaeret-api_default \
 bentsolheim/kilsundvaeret-api:latest
