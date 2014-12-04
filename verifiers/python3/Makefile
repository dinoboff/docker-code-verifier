cwd = $(shell pwd)

repo-name ?= singpath/verifier-python3
tag ?= v0.1

docker ?= docker

images:
	${docker} build --rm=true -t ${repo-name} .
	${docker} build --rm=true -t ${repo-name}:${tag} .

push-images: images
	${docker} login
	${docker} push ${repo-name}:latest
	${docker} push ${repo-name}:${tag}

run-image:
	${docker} run -ti --rm -v ${cwd}:/app ${repo-name}

run-example:
	${docker} run -ti --rm ${repo-name} -e --tests="Pj4+IGZvbwox" "Zm9vPTE="	

test:
	${docker} run -ti --rm -v ${cwd}:/app ${repo-name} python3 -m unittest discover -s codeverifier/

help:
	${docker} run -ti --rm -v ${cwd}:/app ${repo-name} -h