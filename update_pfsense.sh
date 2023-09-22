#!/bin/bash
# This script copies the files from the pfSense/www directory to the pfSense firewall.
# It is intended to be used for development purposes only.

ARGS=$@
PFSENSE_ADDR="192.168.0.1"

if [ -z "$ARGS" ]; then
    echo -e "\nUsage: $0 [USER@]PFSENSE_ADDR\n"
    echo -e "   -h, --help:       Displays this help message"
    echo -e "   -r, --remove:     Removes the pfSense files 1from the pfSense firewall"
    echo -e "   -v, --verbose:    Displays verbose output\n"
    echo -e "   PFSENSE_ADDR:     The hostname or IP of the pfSense firewall to update.\n"
    echo -e "Examples: "
    echo -e "    $0 192.168.0.1"
    echo -e "    $0 root@192.168.0.1"
    exit 1
fi



# Copy files to pfsense
scp -r ./pfSense/www/* root@${PFSENSE_ADDR}:/usr/local/www/