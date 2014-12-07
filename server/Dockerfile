FROM golang:1.3.3-cross

# dependencies
RUN go get github.com/tools/godep && \
	go get code.google.com/p/go.tools/cmd/cover

RUN mkdir -p /go/src/github.com/ChrisBoesch/docker-code-verifier
COPY . /go/src/github.com/ChrisBoesch/docker-code-verifier/server
WORKDIR /go/src/github.com/ChrisBoesch/docker-code-verifier/server

RUN godep restore


CMD ["bash"]