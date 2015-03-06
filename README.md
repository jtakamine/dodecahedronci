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

Etymology
=========
The name "DodecahedronCI" is a reaction against a recent branding trend that favors simple shapes/primitives: Square, Squarespace, Box, Stripe, Line, CircleCI, etc. DodecahedronCI's governing philosophy is to reject the fad of simple lightweight components and favor complicated monolithic systems. (just kidding!)

Setup
=====
Instructions to run DodecahedronCI locally.

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
    Output:
    /Users/myusername/go
    
    $ git clone https://github.com/jtakamine/dodecahedronci.git $GOPATH/src/github.com/jtakamine/dodecahedronci
    
Step 2: Install the client tool
-------------------------------

    $ cd $GOPATH/src/github.com/jtakamine/dodecahedronci/dodec-cli
    $ go get -d && go install
    
Step 3: Start the server
------------------------
The following may take quite a while the first time around. Docker will need to pull relatively large base images if they are not already present in your cache.

    $ cd $GOPATH/src/github.com/jtakamine/dodecahedronci
    $ fig up

If you get an error above, you may need to give your current (non-root) user access to Docker. From the [Docker docs](https://docs.docker.com/installation/ubuntulinux/#giving-non-root-access):

    $ sudo groupadd docker
    $ sudo gpasswd -a ${USER} docker
    
    # If you are in Ubuntu 14.04, replace "docker" with docker.io
    $ sudo service docker restart
    
    $ newgrp docker

    # If you are in Ubuntu 14.04, replace "docker" with "docker.io"
    $ sudo service docker restart
    
Step 4 (boot2docker only): Set CLI target endpoint
--------------------------------------------------
If you are using boot2docker, you will need to override the default dodec-cli target endpoint

    $ export DODEC_ENDPOINT=$(boot2docker ip)
    
Step 5: Trigger a build
-----------------------

    $ dodec-cli execbuild https://github.com/jtakamine/mockrepo.git
    Output:
    UUID
    6ea3ef202b176db5

If you encounter an error above, you may need to add the Go bin directory to your PATH

    $ export PATH=$PATH:$GOPATH/bin

Step 6: View the build logs
---------------------------
In the command below, use the UUID returned in the previous step.

    $ dodec-cli taillogs 6ea3ef202b176db5
    Output:
    2015-03-06T12:25:35Z	 Pulling git repo from https://github.com/jtakamine/mockrepo.git...
    2015-03-06T12:25:35Z	    From https://github.com/jtakamine/mockrepo
    2015-03-06T12:25:35Z	     * branch            HEAD       -> FETCH\_HEAD
    ...

Step 7: Explore
---------------
Check out the [dodec-cli](dodec-cli/) folder for more info about how to use DodecahedronCI.

Architecture
============

The DodecahedronCI server is composed of 6 microservices:
* [dodeccontrol](dodeccontrol/) controls dodecbuild/dodecdeploy, aggregates logs, and exposes the public API.
* [dodecbuild](dodecbuild/) builds Git repos and produces releases.
* [dodecdeploy](dodecdeploy/) deploys releases.
* [dockerregistry](https://github.com/docker/docker-registry) stores docker images.
* [dodecrepo](dodecrepo/) saves/retrieves app data to/from dodecrepodb.
* [dodecrepodb](dodecrepodb/) persists app data (PostgreSQL DB).

DodecahedronCI comes with a CLI client:
* [dodec-cli](dodec-cli/) wraps dodeccontrol's HTTP API in a friendly command-line interface.

![](arch.png)
