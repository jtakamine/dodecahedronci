FROM golang

# Install build essentials
RUN apt-get update && apt-get install -y \
	build-essential

# Install Redis.
RUN \
	cd /tmp && \
	curl -O http://download.redis.io/redis-stable.tar.gz && \
	tar xvzf redis-stable.tar.gz && \
	cd redis-stable && \
	make && \
	make install && \
	cp -f src/redis-sentinel /usr/local/bin && \
	mkdir -p /etc/redis && \
	cp -f *.conf /etc/redis && \
	rm -rf /tmp/redis-stable* && \
	sed -i 's/^\(bind .*\)$/# \1/' /etc/redis/redis.conf && \
	sed -i 's/^\(daemonize .*\)$/# \1/' /etc/redis/redis.conf && \
	sed -i 's/^\(dir .*\)$/# \1\ndir \/data/' /etc/redis/redis.conf && \
	sed -i 's/^\(logfile .*\)$/# \1/' /etc/redis/redis.conf

# Get External Go packages
RUN go get golang.org/x/net/websocket
RUN go get github.com/garyburd/redigo/redis

# Get source and install application
ADD . /go/src/github.com/jtakamine/dodecahedronci/dodecpubsub
WORKDIR /go/src/github.com/jtakamine/dodecahedronci/dodecpubsub
RUN go install

VOLUME ["/data"]

WORKDIR /data

# RUN APPLICATION
CMD redis-server /etc/redis/redis.conf & dodecpubsub -port 8000

EXPOSE 8000
