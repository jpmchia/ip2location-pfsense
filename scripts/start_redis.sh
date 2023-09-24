#!/bin/bash

if [ ! -d "/opt/redis/data" ]; then
  sudo mkdir -p /opt/redis/data
fi

sudo docker run -v /opt/redis/data/:/data -d --name redis-stack -p 6379:6379 -p 8001:8001 -e REDIS_ARGS="--requirepass password" redis/redis-stack:latest
