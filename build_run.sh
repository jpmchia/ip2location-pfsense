#!/bin/bash

mkdir -p ./bin

# Path: build_run.sh
mkdir -p ./local

# Copy the config file if it doesn't exist
if [ ! -f ./local/config.yaml ]; then
    cp -f ./config.yaml.example ./local/config.yaml
fi

# Copy the SSL files if they don't exist
# cp -f ./contrib/ssl/*.pem ./local/
# cp -f ./contrib/ssl/*.key ./local/
# cp -f ./contrib/ssl/*.crt ./local/

# Build the project
go build -o ./bin/ip2location-pfsense ./backend/ && \
cp ./bin/ip2location-pfsense ./local/ip2location-pfsense && \
cd ./local && \
./ip2location-pfsense service $@
