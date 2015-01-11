# docker-code-verifier

This is the docker based code verifier to verify SingPath problem.


## Components


### Code verifier

The verifiers are simple docker images for verifier server. Each image verifies
user solution for a runtime.

By default a verifier server container (using the default `CMD` instruction) 
should bind to port 5000 (for any host).

It should parse the JSON object request inside a `jsonrequest` querystring
( GET request only) or a formdata parameter (POST request), or in the body 
of a POST request.

In the case the Query string the jsonrequest may be url encoded or 
base64 encoded.

It should support JSONP and CORS AJAX request.

The JSON request must have a `solution` attribute and may have a `tests` 
attribute. If the tests are missing the verifier should simply check the 
solution loads correctly.


#### Implemented verifier

- Python 3: tests python code against doctests.
- Angular JS: test a html document against a protractor scenario.


### The web server

The webserver is a simple nginx proxy server (in charge as well of the 
health check).

It's currently ready for testing. See `server/README.md`.


### GAE manager app

Always on. Start and stop a cluster.

TODO:
- [ ] Serves cluster status of the cluster (its load balancer IP and running status)
  those request can used to keep the cluster up.
- [ ] Should start a cluster of code verifier cluster 
  when it receives a status request and the cluster is not running.
- [ ] Should stop the cluster 
  when no cluster status request has been received recently.


## Deployment

Requirements:

- [Docker](https://docs.docker.com/installation/) (or boot2docker on OS X and windows);
- [gcloud](https://cloud.google.com/sdk/#Quick_Start);
- bash and make (you will need to install [cygwin](http://cygwin.com/) on windows).

Note:
	
	On Linux, you might have to run the commands running sudo.


### Testing an instance deployement.

See server's README.md.


### Deploying a cluster

`scripts/deploy.sh` is a bash script creating an instance template, a target pool,
its healthcheck, a network balancer, an instance group and an auto-scaler for the
instance group. To deploy it:
```
export CLUSTER_VERSION=v1.0
cd server; make push-deploy-images tag=$CLUSTER_VERSION; cd -
./scripts/deploy.sh start
```
Note that you can skip the second line if the images of the various containers
have already been uploaded (while testing an instance) and are up-to-date.

Note also that you should be part of the 
[Singpath organization](https://registry.hub.docker.com/repos/singpath/) 
to publish container images on docker hub. 


To stop it:
```
export CLUSTER_VERSION=v1.0
./scripts/deploy.sh clean
```

TODO:
- [ ] The process is rather slow and need improvement.