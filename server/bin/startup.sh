#!/bin/bash

BIN=/usr/local/bin/verifier-server
INIT=/etc/init.d/verifier-server-d
DAEMONUSER=verifier-server
ARCHIVE=server.zip
SRC=./verifier-server

if [[ -z "$CLUSTER_VERSION" ]]; then
    CLUSTER_VERSION=$(curl "http://metadata/computeMetadata/v1/instance/attributes/cluster-version" -H "Metadata-Flavor: Google")
fi

if [[ -z "$CLUSTER_VERSION" ]]; then
    echo "Cluster version is missing. Using 'latest' as default"
    CLUSTER_VERSION="dev"
else
    echo "Cluster version: $CLUSTER_VERSION"
fi

# Istall dependencies
sudo apt-get update
sudo apt-get install -y libcap2-bin make unzip

# pull images
sudo docker pull singpath/verifier-python3

# create user
sudo useradd -s /usr/sbin/nologin -r -M "$DAEMONUSER" > /dev/null 2>&1
sudo adduser "$DAEMONUSER" docker > /dev/null 2>&1

# Install/update server

rm -f "$ARCHIVE"
wget "http://storage.googleapis.com/verifier/server-${CLUSTER_VERSION}" -O "$ARCHIVE"

# start current daemon incase the instance has been rebooted
if [[ -f "$INIT" ]]; then
    sudo "$INIT" stop
fi

if [[ -f "${SRC}/Makefile" ]]; then
    cd "$SRC"; sudo make clean; cd -
fi

# install new version
rm -rf "$SRC"
unzip "$ARCHIVE"
cd verifier-server; sudo make install; cd -

# Start daemon
sudo "$INIT" start
