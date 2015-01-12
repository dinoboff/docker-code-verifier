cwd = $(shell pwd)

repo-name ?= singpath/verifier-angularjs
static-repo-name ?= singpath/verifier-angularjs-static
phantomjs-repo-name ?= singpath/verifier-angularjs-phantomjs
tag ?= latest

docker ?= docker

.images/${repo-name}.${tag}: *.js Dockerfile .dockerignore test/**/*
	mkdir -p .images/singpath
	${docker} build -t ${repo-name}:${tag} .
	touch $@

.images/${static-repo-name}.${tag}: static/*
	mkdir -p .images/singpath
	${docker} build -t ${static-repo-name}:${tag} ./static/
	touch $@

.images/${phantomjs-repo-name}.${tag}: phantomjs/*
	mkdir -p .images/singpath
	${docker} build -t ${phantomjs-repo-name}:${tag} ./phantomjs/
	touch $@

.images/${repo-name}: *.js Dockerfile .dockerignore test/**/*
	mkdir -p .images/singpath
	${docker} build -t ${repo-name} .
	touch $@

.images/${static-repo-name}: static/*
	mkdir -p .images/singpath
	${docker} build -t ${static-repo-name} ./static/
	touch $@

.images/${phantomjs-repo-name}: phantomjs/*
	mkdir -p .images/singpath
	${docker} build -t ${phantomjs-repo-name} ./phantomjs/
	touch $@

images: .images/${repo-name}.${tag} .images/${static-repo-name}.${tag} .images/${phantomjs-repo-name}.${tag} .images/${repo-name} .images/${static-repo-name} .images/${phantomjs-repo-name}

push-images: images
	${docker} login
	${docker} push ${repo-name}:latest
	${docker} push ${repo-name}:${tag}
	${docker} push ${static-repo-name}:latest
	${docker} push ${static-repo-name}:${tag}
	${docker} push ${phantomjs-repo-name}:latest
	${docker} push ${phantomjs-repo-name}:${tag}
.PHONY: images push-images

run-image: images
	${docker} rm -f angularjs-static angularjs-selenium angularjs-phantomjs angularjs || echo "some container were not running..."
	${docker} run -d --name angularjs-static -v /www/_protractor ${static-repo-name}:${tag}
	${docker} run -d --name angularjs-selenium -p 4444:4444 --link angularjs-static:static selenium/hub:2.44.0
	${docker} run -d --name angularjs-phantomjs -h 0.phantomjs.local --link angularjs-selenium:hub --link angularjs-static:static ${phantomjs-repo-name}:${tag}
	${docker} run -ti --rm -p 5000:5000 --link angularjs-selenium:selenium --link angularjs-static:static --volumes-from angularjs-static -v ${cwd}:/app ${repo-name}:${tag} bash

test: images
	${docker} run -ti --rm -v ${cwd}:/app ${repo-name} npm test
.PHONY: run-image test