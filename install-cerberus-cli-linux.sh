#!/bin/sh

echo 'This scripts downloads and installs the latest Cerberus API CLI'
echo 'This script requires curl and jq'

# Fetch the url for the latest linux version
ASSET_URL=$(curl -s https://api.github.com/repos/Nike-Inc/cerberus-cli/releases/latest | \
    jq -r '.assets[] | select(.name=="cerberus-cli-linux-amd64") | .browser_download_url')

echo "Found latest release at ${ASSET_URL}, downloading ..."

# Download the CLI
curl --silent --location --output /usr/local/bin/cerberus ${ASSET_URL}

# Make sure that it is executable
chmod +x /usr/local/bin/cerberus

echo "The Cerberus CLI has been installed to '/usr/local/bin/cerberus', ensure that /usr/local/bin is on your path"