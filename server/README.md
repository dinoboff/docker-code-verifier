# SingPath code verifier server

Run user submitted code in container and return the result.

TODO:

- [ ] Better documentation (including development nitrous).
- [ ] More request type supported. Only support POST Ajax request. It support CORS,
  and the can be send from any domain. The server should support GET 
  and JSONP request.
- [ ] Proper deployment process.


## 1. Go environment

### Local environment

	Note: you can skip this step if you don't need to manage dependencies.
	Using Docker, will allow you to edit, test and build the server.

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

### On Docker

Cross compilation and test are run in a docker container. Although, using a local
Go install might be easier (to use go tools like coverage) or necessary for some tasks
(updating dependencies), you can test and compile the server in docker 
while editing the source locally.

To install the Docker Go environment:

- [install docker](https://docs.docker.com/installation/) (or boot2docker on win/osx).
- `git clone git@github.com:ChrisBoesch/docker-code-verifier.git`.
- optional, pull the Golang image in advance: `docker pull golang:1.3.3-cross`.


	Note: that the [cross compiler Go docker image](https://registry.hub.docker.com/_/golang/) 
	is rather large (1.4 Go). Pulling it the first time starting a Go container might 
	take some time, but running the tests and compile is a rather fast.


## Nitrous environment

TODO.

## 2. Compilation

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

## 3. Testing

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

## 4. Deployment 

```
cd docker-code-verifier/server
make deploy
```

It will:
- test and compile the server executable
- uploaded it to Google Cloud storage
- start a container (using `docker-code-verifier/scripts/deploy.sh`).
- the startup script (`/server/bin/startup.sh`) will pull the verifier images,
  download the server executable and run it (it will be bound to port 80).

You can check if the setup is ready by connecting to the server and 
follow the startup logs (change):
```
gcloud compute --project "singpath-hd" ssh --zone "us-central1-b" "test-verifier"
tail -f /var/log/startupscript.log 
```

Once you see something like:

	Dec 10 00:57:26 test-verifier startupscript: 2014/12/10 00:57:26 Starting server...
	Dec 10 00:57:26 test-verifier startupscript: 2014/12/10 00:57:26 Docker address: unix:///var/run/docker.sock
	Dec 10 00:57:26 test-verifier startupscript: 2014/12/10 00:57:26 Docker cert. path: 
	Dec 10 00:57:26 test-verifier startupscript: 2014/12/10 00:57:26 Binding server to: 0.0.0.0:80


... The server is ready.

To test, use the REST client like Postman and try:

	POST /python HTTP/1.1
	Host: 146.148.37.137
	Cache-Control: no-cache

	{ "solution": "foo=2\nprint(foo)", "tests": ">>> foo\n2" }


The respond should be:

	{
	    "Solved": true,
	    "Printed": "2\n",
	    "Errors": "",
	    "Results": [
	        {
	            "Call": "foo",
	            "Expected": "2",
	            "Received": "2",
	            "Correct": true
	        }
	    ]
	}

The response should be around 400-500 ms .


	Note: the current deployment is only suitable for testing. it will need something
	to monitor the process and restarted it if needed.

## 5. TODO:

- [ ] e2e tests.
- [ ] load testing and tuning concurrent request (set to 5 right now)
- [ ] Proper installation of the server with something like supervisor
  monitoring the process.
- [ ] Create instance template, instance group manager and 
  load balancer for the verifier.
- [ ] Manage startup and shutdown of the cluster via a Google App Engine
  monitor app.
