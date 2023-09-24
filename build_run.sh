#!/bin/bash

mkdir -p ./bin

# Path: build_run.sh
mkdir -p ./local
cp -f ./config.yaml ./local/
cp -f ./counters.yaml ./local/
cp -f ./contrib/ssl/*.pem ./local/
cp -f ./contrib/ssl/*.key ./local/
cp -f ./contrib/ssl/*.crt ./local/

# Build the project
go build -o ./bin/ip2location-pfsense ./backend/ && \
cp ./bin/ip2location-pfsense ./local/ip2location-pfsense && \
cd ./local && \
./ip2location-pfsense service $@
