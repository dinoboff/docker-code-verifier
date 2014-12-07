# SingPath code verifier server

Run user submitted code in container and return the result.


## Go environment

### On Docker

Cross compilation and test are run in a docker container. Although, using a local
Go install might be more easier (to use go tools like coverage) or necessary for some tasks
(updating dependencies), you can test and compile the server in docker 
while editing the source locally.

To install the Docker Go environment:

- [install docker](https://docs.docker.com/installation/) (or boot2docker on win/osx).
- git clone git@github.com:ChrisBoesch/docker-code-verifier.git.
- optional, pull the Golang image in advance: `docker pull golang:1.3.3-cross`.


Note that the cross compile [Go docker image](https://registry.hub.docker.com/_/golang/) 
is rather large (1.4 Go). Pulling it the first starting a Go container might 
take some time, but running the tests and compile is a rather fast.


## Nitrous environment

TODO.


## Local environment

- [install Go](http://golang.org/doc/install) (installed on OS X).
- [install mercurial](http://mercurial.selenic.com/downloads).
- [install git](http://git-scm.com/downloads) (installed on OS X).

You might need SVN too (installed on OS X) (TODO: check if required).

Finally setup [your Go workspace](https://golang.org/doc/code.html#Organization).
Make sure $GOPATH is properly set:
```
echo $GOPATH
```

And clone the repository in $GOPATH:
```
mkdir -p $GOPATH/src/github.com/ChrisBoesch/
cd $GOPATH/src/github.com/ChrisBoesch/
git clone git@github.com:ChrisBoesch/docker-code-verifier.git
go get github.com/tools/godep
cd docker-code-verifier/server
godep restore ./...
```

Note you should still create a docker go environment for cross compiling

## Compile

We will build binary for windows (386/amd64), linux (amd64) and osx (amd64).

We will create a docker container, share the current directory with the container 
and start the container to compile the 4 binaries inside the shared directory. 

```
cd docker-code-verifier/server
make
```

It should have created the 4 binaries:
```
$ ls -la bin/
total 54392
drwxr-xr-x   6 damien  staff   204B  8 Dec 18:15 ./
drwxr-xr-x  10 damien  staff   340B  8 Dec 18:21 ../
-rwxr-xr-x   1 damien  staff   6.9M  8 Dec 18:12 server-linux-amd64
-rwxr-xr-x   1 damien  staff   7.0M  8 Dec 18:12 server-osx-amd64
-rwxr-xr-x   1 damien  staff   5.6M  8 Dec 18:12 server-windows-386.exe
-rwxr-xr-x   1 damien  staff   7.0M  8 Dec 18:12 server-windows-amd64.exe
```

`bin/server-linux-amd64` is the executable we will deploy to the server. The other 
ones are meant for local testing.

## Testing

You can test the libraries in the container:
```
cd docker-code-verifier/server
make test
```

or locally:
```
cd docker-code-verifier/server
go test ./...
```

## Deployment 

TODO:
- post the go server binary somewhere
- start a cluster.
- each node, based of GCE docker VM, should download the server, start it and
 pull the verifier images.

To speed up the creation process of a node we could host the verifier image 
on Google Cloud Storage instead of Docker Hub Registry.

