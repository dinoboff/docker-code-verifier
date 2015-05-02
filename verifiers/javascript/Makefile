cwd = $(shell pwd)

repo-name ?= singpath/verifier-javascript
tag ?= latest
version ?= $(shell cat ../../VERSION)

docker ?= docker

.images/${repo-name}.${tag}: .images/${repo-name}.${version}
	mkdir -p .images/singpath
	${docker} rmi ${repo-name}:${tag} > /dev/null 2>&1 || echo "latest tag not set yet."
	${docker} tag ${repo-name}:${version} ${repo-name}:${tag}
	touch $@

.images/${repo-name}.${version}: Dockerfile *.js test/*.js .dockerignore
	mkdir -p .images/singpath
	${docker} build --rm=true -t ${repo-name}:${version} .
	touch $@

images: .images/${repo-name}.${tag}

push-images: images
	${docker} push ${repo-name}:${version}

.PHONY: images push-images

run-image: images
	${docker} run -ti --rm -p 5000:5000 ${repo-name}:${version} bash

run-image-dev: images
	${docker} run -ti --rm -p 5000:5000 -v ${cwd}:/app ${repo-name}:${version} bash

test: images
	${docker} run -ti --rm ${repo-name}:${version} npm test
.PHONY: run-image test
