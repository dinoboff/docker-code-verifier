#!/bin/bash

VERSION="latest"
SELENIUM_VERSION="2.44.0"

SERVER_IMAGE="singpath/verifier-server"
PYTHON_IMAGE="singpath/verifier-python3"
ANGULARJS_IMAGE="singpath/verifier-angularjs"
ANGULARJS_STATIC_IMAGE="singpath/verifier-angularjs-static"
SELENIUM_IMAGE="selenium/hub"
SELENIUM_PHANTOMJS_IMAGE="singpath/verifier-angularjs-phantomjs"

SERVER_CONTAINER="server"
PYTHON_CONTAINER="python"
ANGULARJS_CONTAINER="angularjs"
ANGULARJS_STATIC_CONTAINER="angularjs-static"
SELENIUM_CONTAINER="selenium"
SELENIUM_PHANTOMJS_CONTAINER="selenium-phantomjs"


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


# Stop packet forward from containers to outside network
# reset
sudo iptables -D FORWARD -i docker0 ! -o docker0 -j ACCEPT
sudo iptables -D FORWARD -i docker0 ! -o docker0 -j DROP
sudo iptables -D FORWARD -i docker0 ! -o docker0 -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
# set rules
sudo iptables -A FORWARD -i docker0 ! -o docker0 -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
sudo iptables -A FORWARD -i docker0 ! -o docker0 -j DROP


# pull images
sudo docker pull "$PYTHON_IMAGE":"$VERSION"
sudo docker pull "$SERVER_IMAGE":"$VERSION"
sudo docker pull "$ANGULARJS_IMAGE":"$VERSION"
sudo docker pull "$ANGULARJS_STATIC_IMAGE":"$VERSION"
sudo docker pull "$SELENIUM_IMAGE":"$SELENIUM_VERSION"
sudo docker pull "$SELENIUM_PHANTOMJS_IMAGE":"$VERSION"


# start verifier containers
sudo docker run -d --name "$PYTHON_CONTAINER" --restart="always"  "$PYTHON_IMAGE":"$VERSION"
sudo docker run -d --name "$ANGULARJS_STATIC_CONTAINER" --restart="always" -v /www/_protractor "$ANGULARJS_STATIC_IMAGE":"$VERSION"
sudo docker run -d --name "$SELENIUM_CONTAINER" --restart="always" --link "$ANGULARJS_STATIC_CONTAINER":static "$SELENIUM_IMAGE":"$SELENIUM_VERSION"
sudo docker run -d --name "$SELENIUM_PHANTOMJS_CONTAINER" --restart="always" -h 0.phantomjs.local --link "$SELENIUM_CONTAINER":hub --link "$ANGULARJS_STATIC_CONTAINER":static "$SELENIUM_PHANTOMJS_IMAGE":"$VERION"
sudo docker run -d --name "$ANGULARJS_CONTAINER" --restart="always" --link "$SELENIUM_CONTAINER":selenium --volumes-from "$ANGULARJS_STATIC_CONTAINER" "$ANGULARJS_IMAGE":"$VERSION"


# start the nginx proxy
sudo docker run -d --name "$SERVER_CONTAINER" -p 80:80 --restart="always" --link "$PYTHON_CONTAINER":python --link "$ANGULARJS_CONTAINER":angularjs "$SERVER_IMAGE":"$VERSION"
