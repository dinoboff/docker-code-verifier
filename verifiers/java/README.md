# Java verifier

The Java verifier needs a Java JDK docker image to compile the server (we use
 `java:8-jdk`) and an image based off a Java JRE (`java:8-jre`).

Those images uses OpenJDK 8.


## Requirements

- [Docker](https://docs.docker.com/installation/) (or docker-machine on OS X
and windows);
- [gcloud](https://cloud.google.com/sdk/#Quick_Start);
- bash and make (you will need to install something [cygwin](http://cygwin.com/)
or [MSYS](http://www.mingw.org/wiki/MSYS) on windows).

Note:

	On Linux, you might have to run the commands running sudo or add your user
	to the docker group.


## Building the docker images

```
make images
```


## Test

The tests are run with docker.

```
make test
```


## Pushing image

You should be part of the
[Singpath organization](https://registry.hub.docker.com/repos/singpath/)
to publish the code verifier on docker hub.

```
make push-images
```