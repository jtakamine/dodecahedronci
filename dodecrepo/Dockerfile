FROM golang

# Get External go packages
RUN go get github.com/lib/pq

# Get Internal go packages
RUN go get github.com/jtakamine/dodecahedronci/logutil #15.03.03.15

ADD . /tmp/

# Get source and install application
ADD . /go/src/github.com/jtakamine/dodecahedronci/dodecrepo
RUN go install github.com/jtakamine/dodecahedronci/dodecrepo

EXPOSE 9000

# Run application
CMD dodecrepo -port 9000
