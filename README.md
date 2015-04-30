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

`scripts/deploy.sh` is a bash script used to setup and start a cluster:
```
$ ./scripts/deploy.sh help
usage: setup|start|stop|delete

setup     Create the an image, an instance template and a load balancer.
start     Manually start the cluster (an instance group and an autoscaler).
          The cluster should already be setup.
stop      Stop the instance group.
delete    Delete the cluster (instance group, load balancer and image.)
          The cluster shouldn't be running.
```

To setup a cluster that Singpath can start:

1. Bump the version in `./VERSION`.

2. Publish the container images:
   ```
   cd server
   make push-deploy-images
   cd ..
   ```

3. Setup the cluster:
   ```
   .scripts/deploy.sh setup
   ```

4. Try the new cluster:
   ```
   .scripts/deploy.sh start
   ```

5. Stop the cluster:
   ```
   .scripts/deploy.sh stop
   ```

An instance template and a network load balancer should now be ready. You just to
update Singpath verifier version. The next time Singpath starts the cluster it
will use the new version.
