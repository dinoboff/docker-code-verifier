cwd = $(shell pwd)

repo-name ?= singpath/verifier-python3
tag ?= latest
version ?= $(shell cat ../../VERSION)

docker ?= docker

.images/${repo-name}.${tag}: .images/${repo-name}.${version}
	mkdir -p .images/singpath
	${docker} rmi ${repo-name}:${tag} > /dev/null 2>&1 || echo "latest tag not set yet."
	${docker} tag -f ${repo-name}:${version} ${repo-name}:${tag}
	touch $@

.images/${repo-name}.${version}: Dockerfile server.py codeverifier/**/*
	mkdir -p .images/singpath
	${docker} build --rm=true -t ${repo-name}:${version} .
	touch $@

images: .images/${repo-name}.${tag}

push-images: images
	${docker} push ${repo-name}:${version}

.PHONY: images push-images

run-image:
	${docker} run -ti --rm -p 5000:5000 -v ${cwd}:/app ${repo-name}:${version}

test:
	${docker} run -ti --rm -v ${cwd}:/app ${repo-name}:${version} python3 -m unittest discover -s codeverifier/
.PHONY: run-image test
