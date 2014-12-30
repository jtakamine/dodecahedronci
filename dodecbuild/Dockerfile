#repoHint:dodecbuild

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

# Get External Go packages
RUN go get gopkg.in/yaml.v2

# Get Internal Go packages
RUN go get github.com/jtakamine/dodecahedronci/configutil
RUN go get github.com/jtakamine/dodecahedronci/dodecregistry/api

# Get source and install application
ADD . /go/src/github.com/jtakamine/dodecahedronci/dodecbuild
WORKDIR /go/src/github.com/jtakamine/dodecahedronci/dodecbuild
RUN go install

# Run application
CMD /usr/local/bin/wrapdocker & dodecbuild -port 8000

EXPOSE 8000
