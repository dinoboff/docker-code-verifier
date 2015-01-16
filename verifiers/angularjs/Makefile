cwd = $(shell pwd)

repo-name ?= singpath/verifier-angularjs
static-repo-name ?= singpath/verifier-angularjs-static
phantomjs-repo-name ?= singpath/verifier-angularjs-phantomjs
tag ?= latest
version ?= $(shell cat ../../VERSION)

docker ?= docker

.images/${repo-name}.${version}: *.js Dockerfile .dockerignore test/**/*
	mkdir -p .images/singpath
	${docker} build -t ${repo-name}:${version} .
	touch $@

.images/${static-repo-name}.${version}: static/*
	mkdir -p .images/singpath
	${docker} build -t ${static-repo-name}:${version} ./static/
	touch $@

.images/${phantomjs-repo-name}.${version}: phantomjs/*
	mkdir -p .images/singpath
	${docker} build -t ${phantomjs-repo-name}:${version} ./phantomjs/
	touch $@

.images/${repo-name}.${tag}: .images/${repo-name}.${version}
	mkdir -p .images/singpath
	${docker} rmi ${repo-name}:${tag} > /dev/null 2>&1 || echo "latest tag not set yet."
	${docker} tag ${repo-name}:${version} ${repo-name}:${tag}
	touch $@

.images/${static-repo-name}.${tag}: .images/${static-repo-name}.${version}
	mkdir -p .images/singpath
	${docker} rmi ${static-repo-name}:${tag} > /dev/null 2>&1 || echo "latest tag not set yet."
	${docker} tag ${static-repo-name}:${version} ${static-repo-name}:${tag}
	touch $@

.images/${phantomjs-repo-name}.${tag}: .images/${phantomjs-repo-name}.${version}
	mkdir -p .images/singpath
	${docker} rmi ${phantomjs-repo-name}:${tag} > /dev/null 2>&1 || echo "latest tag not set yet."
	${docker} tag ${phantomjs-repo-name}:${version} ${phantomjs-repo-name}:${tag}
	touch $@

images: .images/${repo-name}.${tag} .images/${static-repo-name}.${tag} .images/${phantomjs-repo-name}.${tag}

push-images: images
	${docker} push ${repo-name}:${version}
	${docker} push ${static-repo-name}:${version}
	${docker} push ${phantomjs-repo-name}:${version}
.PHONY: images push-images

run-image: images
	${docker} rm -f angularjs-static angularjs-selenium angularjs-phantomjs angularjs || echo "some container were not running..."
	${docker} run -d --name angularjs-static -v /www/_protractor ${static-repo-name}:${version}
	${docker} run -d --name angularjs-selenium -p 4444:4444 --link angularjs-static:static selenium/hub:2.44.0
	${docker} run -d --name angularjs-phantomjs -h 0.phantomjs.local \
			--link angularjs-selenium:hub --link angularjs-static:static ${phantomjs-repo-name}:${version}
	${docker} run -ti --rm -p 5000:5000 --link angularjs-selenium:selenium --link angularjs-static:static \
			--volumes-from angularjs-static -v ${cwd}:/app ${repo-name}:${version} bash

test: images
	${docker} run -ti --rm -v ${cwd}:/app ${repo-name} npm test
.PHONY: run-image test
