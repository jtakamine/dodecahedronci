FROM golang 

# Adapted from https://github.com/tehranian/dind-jenkins-slave/blob/master/Dockerfile
RUN apt-get update -qq && apt-get install -qqy \
	apt-transport-https \
	ca-certificates \
	lxc \
        iptables \
	apparmor
                     
# Install Docker from Docker Inc. repositories.
RUN echo deb https://get.docker.io/ubuntu docker main > /etc/apt/sources.list.d/docker.list \
	&& apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9 \
        && apt-get update -qq \
        && apt-get install -qqy lxc-docker

# Install the magic wrapper
ADD ./utils/wrapdocker /usr/local/bin/wrapdocker
RUN chmod +x /usr/local/bin/wrapdocker

VOLUME /var/lib/docker

# Install stdroute 
RUN go get github.com/jtakamine/stdroute #15.02.16
RUN go install github.com/jtakamine/stdroute

# Get External Go packages
RUN go get gopkg.in/yaml.v2

# Get Internal Go packages
RUN go get github.com/jtakamine/dodecahedronci/logutil #15.02.13.21

# Get source and install application
ADD . /go/src/github.com/jtakamine/dodecahedronci/dodecbuild
RUN go install github.com/jtakamine/dodecahedronci/dodecbuild

EXPOSE 8000

# Run application
CMD /usr/local/bin/wrapdocker & dodecbuild -port 9000 | stdroute -dest $STDROUTE_DEST -method "Logger.Write"
