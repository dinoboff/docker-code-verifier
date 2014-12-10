#!/bin/bash


VM_IMAGE="container-vm-v20141208"
VM_IMAGE_PROJECT="google-containers"
VM_MACHINE_TYPE="f1-micro"

### Authenticaton 

while [[ -z DOCKER_HUB_ACCOUNT ]]; do
	read -p "Your docker user name? " DOCKER_HUB_ACCOUNT
done

echo -e "\n\nChecking your credentials...\n"

CURRENT_USER=$(gcloud config list account | awk '$2 == "=" {print $3}')
if [[ -n "$CURRENT_USER" ]]; then
	echo "You are currently logged it on: $CURRENT_USER"
	read -p "Would you like to login on a different account (y/N)? " yN
	case $yN in
        [Yy]* ) gcloud auth login; ;;
        * ) ;;
    esac
else
	gcloud auth login
fi

### Instance name

while [[ -z $INSTANCE_NAME ]]; do
	echo -e "\n\nSelect an instance name"
	read -p "Instance name? [test-verifier] " INSTANCE_NAME
	if [[ -z $INSTANCE_NAME ]]; then
		INSTANCE_NAME="test-verifier"
	fi
done


while [[ -z "$STARTUP_SCRIPT" ]] && [[ -f "$STARTUP_SCRIPT" ]]; do
	read -p "Startup script for the instance? " STARTUP_SCRIPT
done


echo "creating instance..."

gcloud compute instances create "$INSTANCE_NAME" \
	--image "$VM_IMAGE" \
	--image-project "$VM_IMAGE_PROJECT" \
	--machine-type "$VM_MACHINE_TYPE" \
	--metadata-from-file startup-script="$STARTUP_SCRIPT" \
	--tags http-server
