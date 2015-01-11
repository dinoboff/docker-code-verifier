# SingPath code verifier server

The server is simple proxy. Each runtime verifier runs in a docker container
and a nginx container linked to all verifier container proxy the verification
request to them.


## Requirements

- [Docker](https://docs.docker.com/installation/) (or boot2docker on OS X and windows);
- [gcloud](https://cloud.google.com/sdk/#Quick_Start);
- bash and make (you will need to install [cygwin](http://cygwin.com/) on windows).

Note:
  
  On Linux, you might have to run the commands running sudo.


## Adding a verifier

1. edit `server/bin/startup.sh`:
  
  - startup.sh should pull the  new verifier docker image
  - it to start a new verifier container for the new verifier (each verifier 
    container should have a unique name)
  - it should link the new verifier docker container 
    (add `--link container-name:container-alias` to the argument used to 
    the nginx proxy server).

2. edit server/Makefile various task (`run-image` and `push-deploy-images`) 
   to build, run and publish the new images when testing the deployment of 
   an instance.

3. edit `server/nginx.conf` to add a new `location` like:
   ```location /new-verifier-endpoint {
          proxy_pass    http://container-alias:5000;
   }```
   The new verifier endpoint should be unique. `container-alias` should be
   the one the one set in `server/bin/startup.sh` to link the verifier container
   container to the nginx proxy server container.

4. Upload the server and verifier docker image for the version:
   ```make push-deploy-images```.


## deployment

To test the deployment of a simple instance:
```cd server
export CLUSTER_VERSION=v1.0
make push-deploy-images tag=$CLUSTER_VERSION
make test-deploy tag=CLUSTER_VERSION```

Note that you can skip the third line if the images of the various containers
have already been uploaded (while testing an instance) and are up-to-date.

Note also that you should be part of the 
[Singpath organization](https://registry.hub.docker.com/repos/singpath/) 
to publish container images on docker hub. 


Give a few seconds for the instance to start up, then connect to it over ssh:
```gcloud compute --project <project-id> ssh --zone "us-central1-a" "test-verifier-instance"```

Check the startup logs:
```tail -n +0 -f /var/log/startupscript.log```

Once you see "Finished running startup script /var/run/google.startup.script", 
the server should be running. You can leave the tail (ctrl+c) and end 
the ssh session:
```[ctrl+c]
exit```


You should see some thing like:

  logout
  Connection to 1.2.3.4 closed.

Note the IP address (here, 1.2.3.4) and visit it:
```open http://1.2.3.4```

On OS X it would open a browser at `http://1.2.3.4`; it should display 
"Serving...". Visit `/console/`; you should found a form to test the 
verifiers.
