# docker-code-verifier

This is the docker based code verifier to verify SingPath problem.


## Components


### Code verifier

The verifiers are simple docker images with a common entry point behaviour. 
An entry point should process a payload and return JSON formatted results.

The entry point signature is:

```
docker run --ti verifer-image [-e] [--tests TESTS] solution
```

- `-e`: the input will be base64 encoded;
- `--tests`: the tests to run;
- `solution`: the solution to test.


#### Implemented verifier

- python3


### The web server

The webserver is currently written in Go. 

Working wth Docker and its REST remote API is easy with either Python, Node or Go.
But with the I/O work involved, Nodejs and Go are better candidates. I picked 
Go to speed up deployment.

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
