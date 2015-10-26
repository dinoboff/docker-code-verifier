# Java verifier

The Java verifier needs a Java JDK docker image to compile the server and run
it. We are currently using "java:8-jdk" (OpenJDK).


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

```shell
make images
```


## Running the server

```shell
make run-image
```

To try to POST a JSON solution you can use curl or the Postman chrome extension.
E.g. using curl on OS X (with the docker host IP usually being 192.168.99.100):
```shell
curl \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"solution":"\npublic class SingPath {\n   public Double two() {\n      return 2.0;\n   }\n} \n", "tests":"SingPath sp = new SingPath();\nassertEquals(2.0 , sp.two());"}' \
  http://192.168.99.100:5000/java
```


## Test

The run the tests in docker:

```shell
make test
```


## Pushing image

You should be part of the
[Singpath organization](https://registry.hub.docker.com/repos/singpath/)
to publish the code verifier on docker hub.

```
make push-images
```
