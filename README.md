# docker-code-verifier

This is the docker based code verifier to verify SingPath problem.


## Components


### Code verifier

The verifiers are simple docker images for verifier server. Each image to verifier
user solution for a runtime.

By default a verifier server container (using the default `CMD` instruction) 
should bind to port 5000 (for any host).

 I should be parse the request inside a `jsonrequest` querystring or formdata
 variable or in the json encoded body of a POST request.

 In the case the Query string the jsonrequest may be base64 encoded.

 It should support JSONP and CORS AJAX request.


#### Implemented verifier

- python3


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

TODO