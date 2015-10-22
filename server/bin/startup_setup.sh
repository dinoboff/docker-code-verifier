#!/bin/bash

SELENIUM_VERSION="2.44.0"

SERVER_IMAGE="singpath/verifier-server"
JAVA_IMAGE="singpath/verifier-java"
PYTHON_IMAGE="singpath/verifier-python3"
JAVASCRIPT_IMAGE="singpath/verifier-javascript"
ANGULARJS_IMAGE="singpath/verifier-angularjs"
ANGULARJS_STATIC_IMAGE="singpath/verifier-angularjs-static"
SELENIUM_IMAGE="selenium/hub"
SELENIUM_PHANTOMJS_IMAGE="singpath/verifier-angularjs-phantomjs"


function status_server() {
	sudo mkdir -p /tmp/status
	echo $1 | sudo tee /tmp/status/status.txt
	cd /tmp/status
	sudo python -m SimpleHTTPServer 80
}


# Get version from metadata
if [[ -z "$CLUSTER_VERSION" ]]; then
    CLUSTER_VERSION=$(curl "http://metadata/computeMetadata/v1/instance/attributes/cluster-version" -H "Metadata-Flavor: Google")
fi

if [[ -z "$CLUSTER_VERSION" ]]; then
    echo "Cluster version is missing."
    status_server "failed"
else
	sudo docker pull "$PYTHON_IMAGE":"$CLUSTER_VERSION"
	sudo docker pull "$JAVA_IMAGE":"$CLUSTER_VERSION"
	sudo docker pull "$JAVASCRIPT_IMAGE":"$CLUSTER_VERSION"
	sudo docker pull "$SERVER_IMAGE":"$CLUSTER_VERSION"
	sudo docker pull "$ANGULARJS_IMAGE":"$CLUSTER_VERSION"
	sudo docker pull "$ANGULARJS_STATIC_IMAGE":"$CLUSTER_VERSION"
	sudo docker pull "$SELENIUM_IMAGE":"$SELENIUM_VERSION"
	sudo docker pull "$SELENIUM_PHANTOMJS_IMAGE":"$CLUSTER_VERSION"
	status_server "done"
fi
