cwd = $(shell pwd)

repo-name ?= singpath/verifier-server
tag ?= dev

docker ?= docker

vm-image="container-vm-v20141208"
vm-image-project="google-containers"
vm-machine-type="f1-micro"
vm-zone="us-central1-a"

images:
	${docker} build --rm=true -t ${repo-name} .
	${docker} build --rm=true -t ${repo-name}:${tag} .

push-images: images
	${docker} login
	${docker} push ${repo-name}:latest
	${docker} push ${repo-name}:${tag}

run-image:
	${docker} run -d --name python --restart="always" singpath/verifier-python3:${tag}
	${docker} run -ti --rm -p 8080:80 -v ${cwd}/www:/www --link python:python ${repo-name}
.PHONY: images push-images run-image


push-deploy-images: push-images
	cd ../verifiers/python3; make push-images tag=${tag}

test-deploy: push-deploy-images
	gcloud compute instances create test-verifier-instance \
		--image ${vm-image} \
		--image-project ${vm-image-project} \
		--machine-type ${vm-machine-type} \
		--zone ${vm-zone} \
		--metadata cluster-version="${tag}" \
		--tags http-server
.PHONY: push-deploy-images test-deploy