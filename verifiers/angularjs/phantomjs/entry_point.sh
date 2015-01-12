#!/bin/bash
export GEOMETRY="$SCREEN_WIDTH""x""$SCREEN_HEIGHT""x""$SCREEN_DEPTH"

function shutdown {
	kill -s SIGTERM $NODE_PID
	wait $NODE_PID
}

node_ip=$(cat /etc/hosts | grep phantomjs.local | awk '{print $1}')
node_address=${node_ip}:8080
echo "Phantomjs GhostDriver will bind to " $node_address

started=1
while [[ $started -ne 0 ]]; do
	sleep 1
	echo "checking selenim hub is up..."
	nc -zv hub 4444
	started=$?
done


xvfb-run --server-args="$DISPLAY -screen 0 $GEOMETRY -ac +extension RANDR" \
	phantomjs --webdriver=$node_address --webdriver-selenium-grid-hub=http://hub:4444 
NODE_PID=$!


trap shutdown SIGTERM SIGINT
wait $NODE_PID