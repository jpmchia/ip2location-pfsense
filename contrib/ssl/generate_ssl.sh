#!/bin/bash

# Generate a self-signed certificate for the service to use whilst running on the local host.
#
# This script is intended to be run from the root of the repository.
# It will generate a self-signed certificate for the server to use
# and place it in the correct location.

ROOT_CA_PASSWORD="changeme"

# Generate the root key with the supplied password for the local certificate authority
openssl genrsa -des3 -passout pass:$ROOT_CA_PASSWORD -out localCA.key 4096

# Generate a root-certificate based on the root-key
openssl req -x509 -new -nodes -key localCA.key -passin pass:$ROOT_CA_PASSWORD -config localCA.conf -sha256 -days 365 -out localCA.pem

# Generate a new private key
openssl genrsa -out localhost.key 4096

# Generate a Certificate Signing Request (CSR) based on that private key (reusing the localCA.conf details)
openssl req -new -key local.key -out local.csr -config localCA.conf