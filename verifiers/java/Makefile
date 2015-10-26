cwd = $(shell pwd)

repo-name ?= singpath/verifier-java
tag ?= latest
version ?= $(shell cat ../../VERSION)

docker ?= docker

src-files = $(shell find src -name "*.java")
test-files = $(shell find src -name "*.java")

.images/${repo-name}.${tag}: .images/${repo-name}.${version}
	mkdir -p .images/singpath
	${docker} rmi ${repo-name}:${tag} > /dev/null 2>&1 || echo "latest tag not set yet."
	${docker} tag -f ${repo-name}:${version} ${repo-name}:${tag}
	touch $@

.images/${repo-name}.${version}: Dockerfile lib/*.jar log4j2.xml ${src-files} ${test-files}
	mkdir -p .images/singpath
	${docker} build --rm=true -t ${repo-name}:${version} .
	touch $@

images: .images/${repo-name}.${tag}
.PHONY: images

push-images: images
	${docker} push ${repo-name}:${version}
.PHONY: push-images

run-image: images
	${docker} run -ti --rm -p 5000:5000 ${repo-name}:${version}
.PHONY: run-image

test: .images/${repo-name}.${version}
	$(docker) run -ti --rm ${repo-name}:${version} java -classpath /app/lib/*:/app/build/dist:/app/build/tests org.junit.runner.JUnitCore com.singpath.TestSuite
.PHONY: test

bash: .images/${repo-name}.${version}
	${docker} run -ti --rm -p 5000:5000 ${repo-name}:${version} ls -la /app/src