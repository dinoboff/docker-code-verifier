cwd = $(shell pwd)

repo-name ?= singpath/verifier-python3
tag ?= v0.2

docker ?= docker

images:
	${docker} build --rm=true -t ${repo-name} .
	${docker} build --rm=true -t ${repo-name}:${tag} .

push-images: images
	${docker} login
	${docker} push ${repo-name}:latest
	${docker} push ${repo-name}:${tag}

run-image:
	${docker} run -ti --rm -p 5000:5000 -v ${cwd}:/app ${repo-name}

test:
	${docker} run -ti --rm -v ${cwd}:/app ${repo-name} python3 -m unittest discover -s codeverifier/
