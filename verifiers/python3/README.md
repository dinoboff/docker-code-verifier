# Python3 code verifier

You should have docker [installed](http://docs.docker.com/installation/#installation).

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

You should be part of the singpath organization to publish the code verifier
on docker hub.

```
make push-images
```