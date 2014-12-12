#!/bin/bash

if [[ -z $CLUSTER_VERSION ]]; then
	CLUSTER_VERSION=$(curl "http://metadata/computeMetadata/v1/instance/attributes/cluster-version" -H "Metadata-Flavor: Google")
fi

if [[ -z $CLUSTER_VERSION ]]; then
	echo "Cluster version is missing. Using 'latest' as default"
	CLUSTER_VERSION="latest"
else
	echo "Cluster version: $CLUSTER_VERSION"
fi


sudo docker pull singpath/verifier-python3
wget http://storage.googleapis.com/verifier/server-${CLUSTER_VERSION}
chmod +x server-${CLUSTER_VERSION}
sudo rm -f /usr/local/bin/verifier-server
sudo mv server-${CLUSTER_VERSION} /usr/local/bin/verifier-server
sudo verifier-server -http 0.0.0.0:80
