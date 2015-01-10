cwd = $(shell pwd)

repo-name ?= singpath/verifier-angularjs
static-repo-name ?= singpath/verifier-angularjs-static
tag ?= latest

docker ?= docker

.images/${repo-name}.${tag}: *.js Dockerfile .dockerignore test/**/*
	mkdir -p .images/singpath
	${docker} build --rm=true -t ${repo-name} .
	${docker} build --rm=true -t ${repo-name}:${tag} .
	touch $@

.images/${static-repo-name}.${tag}: static/*
	mkdir -p .images/singpath
	${docker} build -t ${static-repo-name} ./static/
	${docker} build -t ${static-repo-name}:${tag} ./static/
	touch $@

images: .images/${repo-name}.${tag} .images/${static-repo-name}.${tag}

push-images: images
	${docker} login
	${docker} push ${repo-name}:latest
	${docker} push ${repo-name}:${tag}
	${docker} push ${static-repo-name}:latest
	${docker} push ${static-repo-name}:${tag}
.PHONY: images push-images

run-image: images
	${docker} rm -f angularjs-static angularjs-selenium angularjs || echo "some container were not running..."
	${docker} run -d --name angularjs-static -v /www/_protractor ${static-repo-name}:${tag}
	${docker} run -d --name angularjs-selenium --link angularjs-static:static selenium/standalone-chrome:2.44.0
	${docker} run -ti --rm -p 5000:5000 --link angularjs-selenium:selenium --link angularjs-static:static --volumes-from angularjs-static -v ${cwd}:/app ${repo-name}:${tag} bash

test: images
	${docker} run -ti --rm -v ${cwd}:/app ${repo-name} npm test
.PHONY: run-image test