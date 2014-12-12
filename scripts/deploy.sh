#!/bin/bash
#

CLUSTER_VERSION=${VERSION:="dev"}
INSTANCE_GROUP_NAME="verifier-cluster-${CLUSTER_VERSION}"
TEMPLATE_NAME="verifier-template-${CLUSTER_VERSION}"
HEALTHCHECK_NAME="basic-check"
TARGET_POOL_NAME="verifier-pool-${CLUSTER_VERSION}"
AUTOSCALER_NAME="verifier-autoscaler-${CLUSTER_VERSION}"
FORWARD_RULE_NAME="verifier-rule-${CLUSTER_VERSION}"
REGION="us-central1"
ZONE="us-central1-a"

VM_IMAGE="container-vm-v20141208"
VM_IMAGE_PROJECT="google-containers"
VM_MACHINE_TYPE="f1-micro"

CLUSTER_NODE_MIN_COUNT=1
CLUSTER_NODE_MAX_COUNT=2

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
STARTUP_SCRIPT="${DIR}/../server/bin/startup.sh"


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


function group_exist {
    echo -n "Checking if managed instance group named '$1' in zone '$2' already exists..."
    found $(gcloud preview managed-instance-groups --zone "$2" list "$1" | cat -n | awk '$1>1 {print $2}')
    return $?   
}


function healthcheck_exist {
    echo -n "Checking if healthcheck named '$1' already exists..."
    found $(gcloud compute http-health-checks list $1 | cat -n | awk '$1>1 {print $2}')
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


function start() {
    ### Summary
    echo -e "\n\nCluster version: $CLUSTER_VERSION"
    echo "Template name: $TEMPLATE_NAME"
    echo "Image container: $VM_IMAGE"
    echo "Image project: $VM_IMAGE_PROJECT"
    echo "Instance type: $VM_MACHINE_TYPE"
    echo "Number of node: $CLUSTER_NODE_MIN_COUNT"
    echo "Startup script: $STARTUP_SCRIPT"

    ## Authentication
    echo -e "\n\nChecking your credentials...\n"
    set_user

    ### healthcheck
    echo -e "\n\nCreating health-check..."
    create_healthcheck "$HEALTHCHECK_NAME"


    ### Targetpool
    echo -e "\n\nCreating target-pool..."
    create_targetpool "$TARGET_POOL_NAME" "$HEALTHCHECK_NAME" "$REGION"


    ### Instance template
    echo -e "\n\nCreating instance template..."
    create_instance_template "$TEMPLATE_NAME" \
        --machine-type "$VM_MACHINE_TYPE" \
        --image "$VM_IMAGE" \
        --image-project "$VM_IMAGE_PROJECT" \
        --metadata-from-file startup-script="$STARTUP_SCRIPT" \
        --metadata cluster-version="$CLUSTER_VERSION" \
        --tags http-server


    ### Instance group
    echo -e "\n\nCreating instance group..."
    create_group "$INSTANCE_GROUP_NAME" "$ZONE" "$INSTANCE_GROUP_NAME" "$CLUSTER_NODE_MIN_COUNT" "$TEMPLATE_NAME" "$TARGET_POOL_NAME"


    ### Autoscaler
    echo -e "\n\nCreating autoscaler..."
    create_autoscaler "$AUTOSCALER_NAME" "$INSTANCE_GROUP_NAME" "$ZONE" "$CLUSTER_NODE_MIN_COUNT" "$CLUSTER_NODE_MAX_COUNT"


    ### Forwarding rule
    echo -e "\n\nCreating forwarding rule..."
    create_forwardrule "$FORWARD_RULE_NAME" "$REGION" "$TARGET_POOL_NAME"
}


function cleanup() {
    gcloud compute forwarding-rules delete "$TARGET_POOL_NAME" --region "$REGION"
    gcloud compute target-pools delete "$TARGET_POOL_NAME" --region "$REGION"
    gcloud preview autoscaler --zone "$ZONE" delete "$AUTOSCALER_NAME"
    gcloud preview managed-instance-groups --zone "$ZONE" delete "$INSTANCE_GROUP_NAME"
}


function show_help() {
    echo "CMD: start|clean"
}


case "$1" in
    start )
        start
        ;;
    clean )
        cleanup
        ;;
    * )
        show_help
        ;;
esac
