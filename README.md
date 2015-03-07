DodecahedronCI
==============
DodecahedronCI is a continuous integration and deployment server built on Docker and Git.
Through this project, I have introduced myself to the following tools and technologies:
* Go
* Docker/Fig
* PostgreSQL
* JSON-RPC
* Git
* GitHub
* Linux


Quick Start
===========
There is a live DodecahedronCI server running at justintakamine.com. Here's how you can interact with it:

[Install Go](https://golang.org/doc/install) if you don't already have it

Get the dodec-cli client

    $ go get github.com/jtakamine/dodecahedronci/dodec-cli && \
      export DODEC_ENDPOINT=justintakamine.com && \
      export PATH=$PATH:$GOPATH/bin
      
List build history

    $ dodec-cli listbuilds
    
    UUID              APPNAME   VERSION
    4b692f137eaf2ada  mockrepo  0.0.0.1
    c9f59e5b11c4d3fd  mockrepo  0.0.0.2
    ...
    
Execute a new build

    $ dodec-cli execbuild https://github.com/jtakamine/mockrepo.git
    
    UUID
    6ea3ef202b176db5
    
Stream logs from a build

    $ dodec-cli taillogs 6ea3ef202b176db5 #use the UUID returned in command above
    
    2015-03-06T12:25:35Z	 Pulling git repo from https://github.com/jtakamine/mockrepo.git...
    2015-03-06T12:25:35Z	    From https://github.com/jtakamine/mockrepo
    2015-03-06T12:25:35Z	     * branch            HEAD       -> FETCH\_HEAD
    ...
    
Run `dodec-cli help` for more information about how to interact with a DodecahedronCI server.

Etymology
=========
The name "DodecahedronCI" is a reaction against a recent branding trend that favors simple shapes/primitives: Square, Squarespace, Box, Stripe, Line, CircleCI, etc. DodecahedronCI's governing philosophy is to reject the fad of simple lightweight components and favor complicated monolithic systems. (just kidding!)

Architecture
============

The DodecahedronCI server is composed of 6 microservices:
* [dodeccontrol](dodeccontrol/) controls dodecbuild/dodecdeploy, aggregates logs, and exposes the public API.
* [dodecbuild](dodecbuild/) pulls Git repos, builds Dockerfiles, and produces releases.
* [dodecdeploy](dodecdeploy/) deploys releases.
* [dockerregistry](https://github.com/docker/docker-registry) stores docker images. Docker-maintained.
* [dodecrepo](dodecrepo/) saves/retrieves app data to/from dodecrepodb.
* [dodecrepodb](dodecrepodb/) persists app data (PostgreSQL DB).

DodecahedronCI comes with a CLI client:
* [dodec-cli](dodec-cli/) wraps dodeccontrol's HTTP API in a friendly command-line interface.

![](arch.png)



Local Setup
===========
Instructions to run a DodecahedronCI server locally.

Prerequisites
-------------
* [Git](http://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
* [Go](https://golang.org/doc/install)
* [Docker](https://docs.docker.com/installation/)
* [Fig 1.0.1](http://www.fig.sh/install.html)
  * Support for [Compose 1.1.0](https://github.com/docker/compose/releases) TBD

Step 1: Get the source code
---------------------------

    # Make sure the GOPATH environment variable is set before continuing
    $ echo $GOPATH
    
    /Users/myusername/go
    
    $ git clone https://github.com/jtakamine/dodecahedronci.git $GOPATH/src/github.com/jtakamine/dodecahedronci
    
Step 2: Install the CLI client
------------------------------

    $ go get github.com/jtakamine/dodecahedronci/dodec-cli && \
      export PATH=$PATH:$GOPATH/bin
    
    # Only run the following if you are using boot2docker
    $ export DODEC_ENDPOINT=$(boot2docker ip)
    
Step 3: Start the server
------------------------
The following may take quite a while the first time around. Docker will need to pull relatively large base images if they are not already present in your cache.

    $ cd $GOPATH/src/github.com/jtakamine/dodecahedronci && fig up

If you get an error above, you may need to give your current (non-root) user access to Docker. From the [Docker docs](https://docs.docker.com/installation/ubuntulinux/#giving-non-root-access):

    $ sudo groupadd docker
    $ sudo gpasswd -a ${USER} docker
    
    # If you are in Ubuntu, you may need to replace "docker" with "docker.io"
    $ sudo service docker restart
    
    $ newgrp docker
    
Step 4: Trigger a build
-----------------------

    $ dodec-cli execbuild https://github.com/jtakamine/mockrepo.git
    
    UUID
    6ea3ef202b176db5

Step 5: View the build logs
---------------------------
In the command below, use the UUID returned in the previous step.

    $ dodec-cli taillogs 6ea3ef202b176db5
    
    2015-03-06T12:25:35Z	 Pulling git repo from https://github.com/jtakamine/mockrepo.git...
    2015-03-06T12:25:35Z	    From https://github.com/jtakamine/mockrepo
    2015-03-06T12:25:35Z	     * branch            HEAD       -> FETCH\_HEAD
    ...

Step 6: Explore
---------------
Run `dodec-cli help` for more info about how to use DodecahedronCI.

