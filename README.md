# docker-code-verifier

This is the docker based code verifier to verify SingPath problem.


## Components


### Code verifier

The verifiers are simple docker images withe a common entry point behaviour. 
The entry point should process a payload and return JSON formatted results.

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

TODO:
- write a server accepting requests to `POST /<run-time-name>`
  with a json formatted payload like 
  `{"solution": "<solution-code>", "tests": "<tests-code>"}`.


### GAE manager app

Always on. Start and stop a cluster.

TODO:
- Serves cluster status of the cluster (its load balancer IP and running status)
  those request can used to keep the cluster up.
- Should start a cluster of code verifier cluster 
  when it receives a status request and the cluster is not running.
- Should stop the cluster 
  when no cluster status request have been received recently.
