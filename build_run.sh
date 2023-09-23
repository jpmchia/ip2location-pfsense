#!/bin/bash

# Build the project
go build -o ./bin/ip2location-pfsense ./backend/ && ./bin/ip2location-pfsense service $@