# SingPath code verifier server

The server is simple proxy. Each runtime verifier runs in a docker container
and a nginx container linked to all verifier container proxy the verification
request to them.


## Adding a verifier

1. edit `server/bin/startup.sh`:
  
  - startup.sh should pull the  new verifier docker image
  - it to start a new verifier container for the new verifier (each verifier 
    container should have a unique name)
  - it should link the new verifier docker container 
    (add `--link container-name:container-alias` to the argument used to 
    the nginx proxy server).

2. edit `server/Dockerfile` to add a new `location` like:
   ```location /new-verifier-endpoint {
          proxy_pass    http://container-alias:5000;
   }```
   The new verifier endpoint should be unique. `container-alias` should be
   the one the one set in `server/bin/startup.sh` to link the verifier container
   container to the nginx proxy server container.

3. Upload the server and verifier docker image for the version defined in 
   `server/bin/startup.sh` (see the `VERSION` variable).
   You can edit the Makefile rule named `` to do it automatically.


## deploy

TODO:
