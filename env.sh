#!/bin/bash

export $(cat .env)
export TIME_ZONE=Asia/Taipei

# HTTP Configs
export HTTP_LISTEN_ADDR=127.0.0.1
export HTTP_LISTEN_PORT=8081

# Database Configs
export DB_USERNAME=root
export DB_PASSWORD=dorianliu4231
export DB_HOST=127.0.0.1
export DB_PORT=3305
export DB_NAME=kms
