#!/bin/bash

VERSION="dev"
PYTHON_IMAGE="singpath/verifier-python3"
SERVER_IMAGE="singpath/verifier-server"
PYTHON_CONTAINER="python"
SERVER_CONTAINER="server"
DATA=/var/verifier/www

# Get version from metadata
if [[ -z "$CLUSTER_VERSION" ]]; then
    CLUSTER_VERSION=$(curl "http://metadata/computeMetadata/v1/instance/attributes/cluster-version" -H "Metadata-Flavor: Google")
fi

if [[ -z "$CLUSTER_VERSION" ]]; then
    echo "Cluster version is missing. Using default, ${VERSION}."
else
    echo "Cluster version: $CLUSTER_VERSION"
    VERSION="$CLUSTER_VERSION"
fi


#stop packet forward from containers to outside network
sudo iptables -D FORWARD -i docker0 ! -o docker0 -j ACCEPT
sudo iptables -D FORWARD -i docker0 ! -o docker0 -j DROP
sudo iptables -D FORWARD -i docker0 ! -o docker0 -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
sudo iptables -A FORWARD -i docker0 ! -o docker0 -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
sudo iptables -A FORWARD -i docker0 ! -o docker0 -j DROP


# pull images
sudo docker pull "$PYTHON_IMAGE":"$VERSION"
sudo docker pull "$SERVER_IMAGE":"$VERSION"


# Setup server data
sudo mkdir -p "$DATA"
rm -rf "${DATA}/index.html"
echo "<html>serving...</html>" | sudo tee "${DATA}/index.html"


# start verifier containers
sudo docker run -d --name "$PYTHON_CONTAINER" --restart="always"  "$PYTHON_IMAGE":"$VERSION"


# start the nginx proxy
sudo docker run -d --name "$SERVER_CONTAINER" -p 80:80 -v "$DATA":/www --restart="always" --link "$PYTHON_CONTAINER":python "$SERVER_IMAGE":"$VERSION"
