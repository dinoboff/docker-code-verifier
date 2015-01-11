cwd = $(shell pwd)

repo-name ?= singpath/verifier-server
tag ?= latest

docker ?= docker

vm-image="container-vm-v20141208"
vm-image-project="google-containers"
vm-machine-type="f1-micro"
vm-zone="us-central1-a"

.images/${repo-name}: Dockerfile nginx.conf www/**/*
	mkdir -p .images/singpath
	${docker} build --rm=true -t ${repo-name} .
	touch $@

.images/${repo-name}.${tag}: Dockerfile nginx.conf www/**/*
	mkdir -p .images/singpath
	${docker} build --rm=true -t ${repo-name}:${tag} .
	touch $@

images: .images/${repo-name} .images/${repo-name}.${tag}

push-images: images
	${docker} login
	${docker} push ${repo-name}:latest
	${docker} push ${repo-name}:${tag}

run-image:
	cd ../verifiers/angularjs; make images tag=${tag}
	cd ../verifiers/python3; make images tag=${tag}
	${docker} run -d --name python singpath/verifier-python3:${tag}
	${docker} run -d --name angularjs-static -v /www/_protractor singpath/verifier-angularjs-static:${tag}
	${docker} run -d --name angularjs-selenium --link angularjs-static:static selenium/hub:2.44.0
	${docker} run -d --name angularjs-phantomjs -h container.host -p 5555:5555 -e NODE_PORT=5555 --link angularjs-selenium:hub --link angularjs-static:static singpath/verifier-angularjs-phantomjs:${tag}
	${docker} run -d --name angularjs --link angularjs-selenium:selenium --link angularjs-static:static --volumes-from angularjs-static singpath/verifier-angularjs:latest
	${docker} run -ti --rm -p 8080:80 -v ${cwd}/www:/www --link python:python --link angularjs:angularjs ${repo-name}
.PHONY: images push-images run-image

push-deploy-images: push-images
	cd ../verifiers/python3; make push-images tag=${tag}
	cd ../verifiers/angularjs; make push-images tag=${tag}

test-deploy:
	gcloud compute instances create test-verifier-instance \
		--image ${vm-image} \
		--image-project ${vm-image-project} \
		--machine-type ${vm-machine-type} \
		--zone ${vm-zone} \
		--metadata-from-file startup-script="bin/startup.sh" \
		--metadata cluster-version="${tag}" \
		--tags http-server
.PHONY: push-deploy-images test-deploy