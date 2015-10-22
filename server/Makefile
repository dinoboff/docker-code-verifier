cwd = $(shell pwd)

repo-name ?= singpath/verifier-server
tag ?= latest
version ?= $(shell cat ../VERSION)

docker ?= docker

vm-image="container-vm-v20150317"
vm-image-project="google-containers"
vm-machine-type="f1-micro"
vm-zone="us-central1-a"

.images/${repo-name}.${tag}: .images/${repo-name}.${version}
	mkdir -p .images/singpath
	${docker} rmi ${repo-name}:${tag} > /dev/null 2>&1 || echo "latest tag not set yet."
	${docker} tag -f ${repo-name}:${version} ${repo-name}:${tag}
	touch $@

.images/${repo-name}.${version}: Dockerfile nginx.conf www/**/*
	mkdir -p .images/singpath
	${docker} build --rm=true -t ${repo-name}:${version} .
	touch $@

images: .images/${repo-name}.${tag}

push-images: images
	${docker} login
	${docker} push ${repo-name}:${version}

run-image: images
	${docker} rm -f $(${docker} ps -aq) > /dev/null 2>&1 || echo "no container running."
	cd ../verifiers/angularjs; make images version=${version} tag=${tag}
	cd ../verifiers/java; make images version=${version} tag=${tag}
	cd ../verifiers/python3; make images version=${version} tag=${tag}
	cd ../verifiers/javascript; make images version=${version} tag=${tag}
	CLUSTER_VERSION=${version} SKIP_CLUSTER_IPTABLE_CONFIGURATION=true DOCKER="docker" ./bin/startup_run.sh
.PHONY: images push-images run-image

push-deploy-images: push-images
	cd ../verifiers/java; make push-images version=${version} tag=${tag}
	cd ../verifiers/python3; make push-images version=${version} tag=${tag}
	cd ../verifiers/javascript; make push-images version=${version} tag=${tag}
	cd ../verifiers/angularjs; make push-images version=${version} tag=${tag}

test-deploy:
	gcloud compute instances create test-verifier-instance \
		--image ${vm-image} \
		--image-project ${vm-image-project} \
		--machine-type ${vm-machine-type} \
		--zone ${vm-zone} \
		--metadata-from-file startup-script="bin/startup_run.sh" \
		--metadata cluster-version="${version}" \
		--tags http-server
.PHONY: push-deploy-images test-deploy
