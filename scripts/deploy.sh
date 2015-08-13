#!/bin/bash
#

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

VERSION=$(cat "${DIR}/../VERSION")
CLUSTER_VERSION=${CLUSTER_VERSION:=$VERSION}

INSTANCE_BUILDER_NAME="verifier-instance-image-builder-${CLUSTER_VERSION//./-}"
INSTANCE_GROUP_NAME="verifier-cluster-${CLUSTER_VERSION//./-}"
TEMPLATE_NAME="verifier-template-${CLUSTER_VERSION//./-}"
HEALTHCHECK_NAME="basic-check"
TARGET_POOL_NAME="verifier-pool-${CLUSTER_VERSION//./-}"
AUTOSCALER_NAME="verifier-autoscaler-${CLUSTER_VERSION//./-}"
FORWARD_RULE_NAME="verifier-rule-${CLUSTER_VERSION//./-}"
REGION="us-central1"
ZONE="us-central1-a"

BASE_IMAGE="container-vm-v20150806"
BASE_IMAGE_PROJECT="google-containers"
VM_IMAGE="verifier-${CLUSTER_VERSION//\./-}"
VM_IMAGE_PROJECT="singpath-hd"
VM_MACHINE_TYPE="f1-micro"

CLUSTER_NODE_MIN_COUNT=1
CLUSTER_NODE_MAX_COUNT=2

STARTUP_SCRIPT="${DIR}/../server/bin/startup_run.sh"
STARTUP_SETUP_SCRIPT="${DIR}/../server/bin/startup_setup.sh"


### Test Startup script
if [[ -z "$STARTUP_SCRIPT" ]] || [[ ! -f "$STARTUP_SCRIPT" ]]; then
    >&2 echo "The startup script could not be found."
    exit 1
fi


function set_user() {
    current_user=$(gcloud config list account | awk '$2 == "=" {print $3}')
    if [[ -n "$current_user" ]]; then
        echo "You are currently logged it on: $current_user"
        read -p "Would you like to login on a different account (y/N)? " yN
        case $yN in
            [Yy]* ) gcloud auth login; ;;
            * ) ;;
        esac
    else
        gcloud auth login
    fi
}


function found {
    if [[ -z "$1" ]]; then
        echo " not found."
        return 1
    else
        echo " found."
        return 0
    fi
}


function forwardrule_exist() {
    echo -n "Checking if forwarding rule named '$1' in region '$2' already exists..."
    found $(gcloud compute forwarding-rules list "$1" --regions "$2" | cat -n | awk '$1>1 {print $2}')
    return $?
}


function group_exist() {
    echo -n "Checking if managed instance group named '$1' in zone '$2' already exists..."
    found $(gcloud preview managed-instance-groups --zone "$2" list -l | awk '$1=="'$1'" {print $1}')
    return $?   
}


function healthcheck_exist() {
    echo -n "Checking if healthcheck named '$1' already exists..."
    found $(gcloud compute http-health-checks list $1 | cat -n | awk '$1>1 {print $2}')
    return $?
}

function image_exist() {
    echo -n "Checking if image named '$1' already exists..."
    found $(gcloud compute images list --no-standard-images $1 | cat -n | awk '$1>1 {print $2}')
    return $?
}


function targetpool_exist() {
    echo -n "Checking if target pool named '$1' in region $2 already exists..."
    found $(gcloud compute target-pools list "$1" --regions "$2" | cat -n | awk '$1>1 {print $2}')
    return $?
}


function template_exist() {
    echo -n "Checking if template named '$1' already exists..."
    found $(gcloud compute instance-templates list $1 | cat -n | awk '$1>1 {print $2}')
    return $?
}


function create_autoscaler() {
    gcloud preview autoscaler --zone "$3" create "$1" \
         --min-num-replicas "$4" \
         --max-num-replicas "$5" \
         --target "$2"
}


function create_forwardrule() {
    forwardrule_exist "$1" "$2"
    if [[ $? -ne 0 ]]; then
        gcloud compute forwarding-rules create "$1" \
            --region "$2" \
            --port-range 80 \
            --target-pool "$3"
    else
        >&2 echo "The forwarding rule should be removed (is it safe?)"
        >&2 echo "or you could deploy to an other version"
        exit 1
    fi
}


function create_group() {
    group_exist $1 $2
    if [[ $? -ne 0 ]]; then
        gcloud preview managed-instance-groups \
            --zone "$2" \
            create "$1" \
            --base-instance-name "$3" \
            --size  "$4"\
            --template "$5" \
            --target-pool "$6"
    else
        >&2 echo "The instance group exists already."
        exit 1
    fi
}


function create_healthcheck() {
    healthcheck_exist $1
    if [[ $? -ne 0 ]]; then
        gcloud compute http-health-checks create $1
    fi
}


function create_instance_template() {
    template_exist $TEMPLATE_NAME
    if [[ $? -eq 0 ]]; then
        gcloud compute instance-templates delete $1
    fi

    echo -e "\nCreateing instance template..."
    gcloud compute instance-templates create $@
}


function create_targetpool() {
    targetpool_exist $1 $3
    if [[ $? -ne 0 ]]; then
        gcloud compute target-pools create $1 \
            --region $3 --health-check $2
    else
        >&2 echo "The target pool should be removed (is it safe?)"
        >&2 echo "or you could deploy to an other version"
        exit 1
    fi
}


function create_image() {
    image_exist $1
    if [[ $? -ne 0 ]]; then
        start_image_builder_instance $1
        check_image_builder_instance_status $1
        save_image_builder_instance $1
    else
        echo "Image $INSTANCE_BUILDER_NAME already exists."
        echo "You should delete it or prepare a new version."
        exit 1
    fi
}


function start_image_builder_instance() {
    gcloud compute instances create "$1" \
        --image "$BASE_IMAGE" \
        --image-project "$BASE_IMAGE_PROJECT" \
        --machine-type "$VM_MACHINE_TYPE" \
        --zone "$ZONE" \
        --metadata-from-file startup-script="$STARTUP_SETUP_SCRIPT" \
        --metadata cluster-version="$CLUSTER_VERSION" \
        --tags http-server
}


function check_image_builder_instance_status() {
    details=$(gcloud compute instances list --regexp "$1" --zone "$ZONE")
    ip=$(echo "$details" | cat -n | awk '$1>1 {print $6}')
    if [[ -z "$ip" ]]; then
        echo "instance ip not found"
        exit 1
    fi

    echo "Image builder instance IP:" $ip
    
    status=$(curl http://${ip}/status.txt)
    status_fetched=$?
    while [[ $status_fetched -ne 0 ]]; do
        echo "Failed to fetch image builder instance status..."
        echo "will try again in 30s"
        sleep 30
        status=$(curl http://${ip}/status.txt)
        status_fetched=$?
    done

    if [[ "$status" == "done" ]]; then
        echo "Image builder instance is ready."
    elif [[ "$status" == "failed" ]]; then
        echo "Image builder instance startup failed."
        exit 2
    else
        echo "Unknown image builder instance status:" $status
        exit 3
    fi
}


function save_image_builder_instance() {
    echo "Stopping image builder instance..."
    gcloud compute instances delete "$1" --zone "$ZONE" --keep-disks boot

    echo "Creating a new image named $VM_IMAGE from the the boot disk..."
    gcloud compute images create "$VM_IMAGE" --source-disk "$1" --source-disk-zone "$ZONE"

    echo "Deleting the boot disk..."
    gcloud compute disks delete "$1" --zone "$ZONE"
}


function setup_cluster() {
    ### Summary
    echo -e "\n\nCluster version: $CLUSTER_VERSION"
    echo "Template name: $TEMPLATE_NAME"
    echo "Base Image container: $BASE_IMAGE"
    echo "Base Image project: $BASE_IMAGE_PROJECT"
    echo "Instance type: $VM_MACHINE_TYPE"
    echo "Startup script: $STARTUP_SCRIPT"

    ### Authentication
    echo -e "\n\nChecking your credentials...\n"
    set_user

    ### Verifier image, step 1
    echo -e "\n\nStarting creating the verifier image..."
    image_exist "$INSTANCE_BUILDER_NAME"
    if [[ $? -eq 0 ]]; then
        echo "Image $INSTANCE_BUILDER_NAME already exists."
        echo "You should delete it or prepare a new version."
        exit 1
    fi
    start_image_builder_instance "$INSTANCE_BUILDER_NAME"


    ### healthcheck
    echo -e "\n\nCreating health-check..."
    create_healthcheck "$HEALTHCHECK_NAME"

    ### Targetpool
    echo -e "\n\nCreating target-pool..."
    create_targetpool "$TARGET_POOL_NAME" "$HEALTHCHECK_NAME" "$REGION"

    ### Forwarding rule
    echo -e "\n\nCreating forwarding rule..."
    create_forwardrule "$FORWARD_RULE_NAME" "$REGION" "$TARGET_POOL_NAME"

    ### Verifier image, step 2 and 3
    echo -e "\n\nFinishing creating the verifier image..."
    check_image_builder_instance_status "$INSTANCE_BUILDER_NAME"
    save_image_builder_instance "$INSTANCE_BUILDER_NAME"

    ### Instance template
    echo -e "\n\nCreating instance template..."
    create_instance_template "$TEMPLATE_NAME" \
        --machine-type "$VM_MACHINE_TYPE" \
        --image "$VM_IMAGE" \
        --image-project "$VM_IMAGE_PROJECT" \
        --metadata-from-file startup-script="$STARTUP_SCRIPT" \
        --metadata cluster-version="$CLUSTER_VERSION" \
        --tags http-server
}

function start_cluster() {
    ### Instance group
    echo -e "\nCreating instance group..."
    create_group "$INSTANCE_GROUP_NAME" "$ZONE" "$INSTANCE_GROUP_NAME" "$CLUSTER_NODE_MIN_COUNT" "$TEMPLATE_NAME" "$TARGET_POOL_NAME"


    ### Autoscaler
    echo -e "\n\nCreating autoscaler..."
    create_autoscaler "$AUTOSCALER_NAME" "$INSTANCE_GROUP_NAME" "$ZONE" "$CLUSTER_NODE_MIN_COUNT" "$CLUSTER_NODE_MAX_COUNT"

    cluster_ip=$(gcloud compute forwarding-rules list codeverifier-rule --regions us-central1 | cat -n | awk '$1>1{print $4}')
    echo "Test the cluster at http://${cluster_ip}/console/"
    echo "(it may take a minute for the cluster to be available)"
}


function stop_cluster() {
    gcloud preview autoscaler --zone "$ZONE" delete "$AUTOSCALER_NAME"
    gcloud preview managed-instance-groups --zone "$ZONE" delete "$INSTANCE_GROUP_NAME"
}


function delete_cluster() {
    gcloud compute forwarding-rules delete "$FORWARD_RULE_NAME" --region "$REGION"
    gcloud compute target-pools delete "$TARGET_POOL_NAME" --region "$REGION"
    gcloud compute instance-templates delete "$TEMPLATE_NAME"
    gcloud compute images delete "$VM_IMAGE"
}


function show_help() {
    echo -e "usage: setup|start|stop|delete \n"
    echo "setup     Create the an image, an instance template and a load balancer."
    echo "start     Manually start the cluster (an instance group and an autoscaler)."
    echo "          The cluster should already be setup."
    echo "stop      Stop the instance group."
    echo "delete    Delete the cluster (instance group, load balancer and image.)"
    echo -e "          The cluster shouldn't be running.\n"
}


case "$1" in
    setup )
        setup_cluster
        ;;
    start )
        start_cluster
        ;;
    stop )
        stop_cluster
        ;;
    delete )
        delete_cluster
        ;;
    * )
        show_help
        ;;
esac
