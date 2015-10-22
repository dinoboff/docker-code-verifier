#!/bin/bash

DOCKER=${DOCKER:="sudo docker"}
IPTABLES=${IPTABLES:="sudo iptables"}

SELENIUM_VERSION="2.44.0"

SERVER_IMAGE="singpath/verifier-server"
JAVA_IMAGE="singpath/verifier-java"
PYTHON_IMAGE="singpath/verifier-python3"
JAVASCRIPT_IMAGE="singpath/verifier-javascript"
ANGULARJS_IMAGE="singpath/verifier-angularjs"
ANGULARJS_STATIC_IMAGE="singpath/verifier-angularjs-static"
SELENIUM_IMAGE="selenium/hub"
SELENIUM_PHANTOMJS_IMAGE="singpath/verifier-angularjs-phantomjs"

SERVER_CONTAINER="server"
JAVA_CONTAINER="java"
PYTHON_CONTAINER="python"
JAVASCRIPT_CONTAINER="javascript"
ANGULARJS_CONTAINER="angularjs"
ANGULARJS_STATIC_CONTAINER="angularjs-static"
SELENIUM_CONTAINER="selenium"
SELENIUM_PHANTOMJS_CONTAINER="selenium-phantomjs"


# Get version from metadata
if [[ -z "$CLUSTER_VERSION" ]]; then
    CLUSTER_VERSION=$(curl "http://metadata/computeMetadata/v1/instance/attributes/cluster-version" -H "Metadata-Flavor: Google")
fi

if [[ -z "$CLUSTER_VERSION" ]]; then
    echo "Cluster version is missing."
    exit 1
fi

if [[ -z "$SKIP_CLUSTER_IPTABLE_CONFIGURATION" ]]; then
	# Stop packet forward from containers to outside network
	# reset
	$IPTABLES -D FORWARD -i docker0 ! -o docker0 -j ACCEPT
	$IPTABLES -D FORWARD -i docker0 ! -o docker0 -j DROP
	$IPTABLES -D FORWARD -i docker0 ! -o docker0 -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
	# set rules
	$IPTABLES -A FORWARD -i docker0 ! -o docker0 -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
	$IPTABLES -A FORWARD -i docker0 ! -o docker0 -j DROP
fi


# Remove all containers
containers=$($DOCKER ps -aq)
if [[ -n "$containers" ]]; then
	$DOCKER rm -f $containers
fi


# start verifier containers
$DOCKER run -d --name "$JAVA_CONTAINER" --restart="always"  "$JAVA_IMAGE":"$CLUSTER_VERSION"
$DOCKER run -d --name "$PYTHON_CONTAINER" --restart="always"  "$PYTHON_IMAGE":"$CLUSTER_VERSION"
$DOCKER run -d --name "$JAVASCRIPT_CONTAINER" --restart="always"  "$JAVASCRIPT_IMAGE":"$CLUSTER_VERSION"
$DOCKER run -d --name "$ANGULARJS_STATIC_CONTAINER" --restart="always" -v /www/_protractor "$ANGULARJS_STATIC_IMAGE":"$CLUSTER_VERSION"
$DOCKER run -d --name "$SELENIUM_CONTAINER" --restart="always" --link "$ANGULARJS_STATIC_CONTAINER":static "$SELENIUM_IMAGE":"$SELENIUM_VERSION"
$DOCKER run -d --name "$SELENIUM_PHANTOMJS_CONTAINER" --restart="always" -h 0.phantomjs.local --link "$SELENIUM_CONTAINER":hub --link "$ANGULARJS_STATIC_CONTAINER":static "$SELENIUM_PHANTOMJS_IMAGE":"$CLUSTER_VERSION"
$DOCKER run -d --name "$ANGULARJS_CONTAINER" --restart="always" --link "$SELENIUM_CONTAINER":selenium --volumes-from "$ANGULARJS_STATIC_CONTAINER" "$ANGULARJS_IMAGE":"$CLUSTER_VERSION"


# start the nginx proxy
$DOCKER run -d --name "$SERVER_CONTAINER" -p 80:80 --restart="always" \
	--link "$JAVA_CONTAINER":java \
	--link "$PYTHON_CONTAINER":python \
	--link "$JAVASCRIPT_CONTAINER":javascript \
	--link "$ANGULARJS_CONTAINER":angularjs \
	"$SERVER_IMAGE":"$CLUSTER_VERSION"
