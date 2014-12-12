cwd = $(shell pwd)

package-name = github.com/ChrisBoesch/docker-code-verifier/server
repo-name ?= singpath/verifier-server
tag ?= v0.1
gosrc = $(shell find . -name "*.go" | sed 's/ /\\ /g')
target ?= bin/server-osx-amd64
version ?= dev

dist-name = verifier-server
dist-zip = ${dist-name}.zip
dist-path = dist/${dist-zip}

vm-image="container-vm-v20141208"
vm-image-project="google-containers"
vm-machine-type="f1-micro"
vm-zone="us-central1-a"

docker ?= docker

build: images test bin/server-linux-amd64
	rm -rf build
	mkdir -p build
	cp bin/{server-linux-amd64,init.sh,Makefile} build/

dist: build
	rm -rf dist
	mkdir -p dist/${dist-name}
	cp build/* dist/${dist-name}
	cd dist; zip ./${dist-zip} ${dist-name}/*
	rm -rf dist/${dist-name}

dev: images test ${target}
	echo "starting server..."
	${target}

.images/${repo-name}: Dockerfile
	mkdir -p .images/singpath/
	${docker} build --rm=true -t ${repo-name} .
	touch $@

.images/${repo-name}.${tag}: Dockerfile
	mkdir -p .images
	${docker} build --rm=true -t ${repo-name}:${tag} .
	touch $@

images: .images/${repo-name} .images/${repo-name}.${tag}

run-image:
	${docker} run --rm -ti -v ${cwd}:/go/src/${package-name} ${repo-name}

bin/server-windows-386.exe: ${gosrc}
	mkdir -p ./bin
	rm -rf $@
	${docker} run --rm -ti -v ${cwd}:/go/src/${package-name} -e GOOS=windows -e GOARCH=386 ${repo-name} go build -v -o $@

bin/server-windows-amd64.exe: ${gosrc}
	mkdir -p ./bin
	rm -rf $@
	${docker} run --rm -ti -v ${cwd}:/go/src/${package-name} -e GOOS=windows -e GOARCH=amd64 ${repo-name} go build -v -o $@

bin/server-linux-amd64: ${gosrc}
	mkdir -p ./bin
	rm -rf $@
	${docker} run --rm -ti -v ${cwd}:/go/src/${package-name} -e GOOS=linux -e GOARCH=amd64 ${repo-name} go build -v -o $@

bin/server-osx-amd64: ${gosrc}
	mkdir -p ./bin
	rm -rf $@
	${docker} run --rm -ti -v ${cwd}:/go/src/${package-name} -e GOOS=darwin -e GOARCH=amd64 ${repo-name} go build -v -o $@

test:
	${docker} run --rm -ti -v ${cwd}:/go/src/${package-name} ${repo-name} go test ${package-name}/...

upload: dist
	gsutil cp ${dist-path} gs://verifier/server-${version}
	gsutil acl set public-read gs://verifier/server-${version}

deploy-test:
	gcloud compute instances create test-verifier-instance \
		--image ${vm-image} \
		--image-project ${vm-image-project} \
		--machine-type ${vm-machine-type} \
		--zone ${vm-zone} \
		--tags http-server
