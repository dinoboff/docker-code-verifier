J = java
JFLAGS = -classpath lib/*:.
JC = javac
JCFLAGS = $(JFLAGS) -sourcepath ./src -d ./build/dist -g
JCTESTFLAGS = -classpath "lib/*:build/dist" -sourcepath ./tests -d ./build/tests -g
.SUFFIXES:

CLASSES = src/com/singpath/Failure.java \
	src/com/singpath/Request.java \
	src/com/singpath/Server.java \
	src/com/singpath/Verifier.java \
	src/com/singpath/OnShutDown.java \
	src/com/singpath/Response.java \
	src/com/singpath/ShutDown.java \
	src/com/singpath/VerifierHandler.java
TESTS = tests/com/singpath/FailureTest.java \
	tests/com/singpath/RequestTest.java \
	tests/com/singpath/ResponseTest.java \
	tests/com/singpath/TestSuite.java \
	tests/com/singpath/VerifierTest.java

cwd = $(shell pwd)

repo-name ?= singpath/verifier-java
tag ?= latest
version ?= $(shell cat ../../VERSION)

docker ?= docker
DFLAGS = run -it --rm -v ${cwd}:/app
BUILDERIMAGE = java:8-jdk

default: classes

classes: $(CLASSES:src/%.java=build/dist/%.class)
tests-classes: $(TESTS:tests/%.java=build/tests/%.class)
.PHONY: classes tests-classes

build/dist/%.class: src/%.java
	@mkdir -p build/dist
	$(docker) $(DFLAGS) -w /app $(BUILDERIMAGE) $(JC) $(JCFLAGS) src/$*.java

build/tests/%.class: tests/%.java
	@mkdir -p build/tests
	$(docker) $(DFLAGS) -w /app $(BUILDERIMAGE) $(JC) $(JCTESTFLAGS) tests/$*.java

build/dist/verifier.jar: $(CLASSES:src/%.java=build/dist/%.class)
	rm -rf build/dist/verifier.jar
	$(docker) $(DFLAGS) -w /app/build/dist $(BUILDERIMAGE) jar cvfe verifier.jar com.singpath.Server .

bash:
	$(docker) $(DFLAGS) -w /app $(BUILDERIMAGE) bash

build: build/dist/verifier.jar
.PHONY: build

clean:
	rm -rf build
.PHONY: clean

test: classes tests-classes
	$(docker) $(DFLAGS) -w /app $(BUILDERIMAGE) $(J) -classpath lib/*:build/dist/:build/tests/ org.junit.runner.JUnitCore com.singpath.TestSuite
.PHONY: test

.images/${repo-name}.${tag}: .images/${repo-name}.${version}
	mkdir -p .images/singpath
	${docker} rmi ${repo-name}:${tag} > /dev/null 2>&1 || echo "latest tag not set yet."
	${docker} tag -f ${repo-name}:${version} ${repo-name}:${tag}
	touch $@

.images/${repo-name}.${version}: Dockerfile lib/*.jar log4j2.xml build/dist/verifier.jar
	mkdir -p .images/singpath
	${docker} build --rm=true -t ${repo-name}:${version} .
	touch $@

images: .images/${repo-name}.${tag}

push-images: images
	${docker} push ${repo-name}:${version}

.PHONY: images push-images

run-image: images
	${docker} run -ti --rm -p 5000:5000 ${repo-name}:${version}
