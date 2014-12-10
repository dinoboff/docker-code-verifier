#!/bin/bash

sudo docker pull singpath/verifier-python3
wget http://storage.googleapis.com/verifier/server
chmod +x server
sudo rm -f /usr/local/bin/verifier-server
sudo mv server /usr/local/bin/verifier-server
sudo verifier-server -http 0.0.0.0:80
