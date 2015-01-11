# Angular JS verifier

The Angular JS verifier require 3 types on container:

- the verifier server: it receives the json requests, parse them and run the 
  payload via protractor;
- the selenium hub server: protractor will send the tests to the selenium server;
- the selenium nodes: currently only one Phantomjs node.

Note that tests are currently run serially. To run concurrent tests in parallel,
we will need more than one phantomjs container running.


TODO:
- [ ] check how many Phantomjs container can run on the small GCE instances.


## Requirements

- [Docker](https://docs.docker.com/installation/) (or boot2docker on OS X and windows);
- [gcloud](https://cloud.google.com/sdk/#Quick_Start);
- bash and make (you will need to install [cygwin](http://cygwin.com/) on windows).

Note:
	
	On Linux, you might have to run the commands running sudo.


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