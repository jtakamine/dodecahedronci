FROM golang

# Install stdroute
RUN go get github.com/jtakamine/stdroute
RUN go install github.com/jtakamine/stdroute

# Get External go packages
RUN go get github.com/lib/pq

# Get Internal go packages
RUN go get github.com/jtakamine/dodecahedronci/logutil

# Get source and install application
ADD . /go/src/github.com/jtakamine/dodecahedronci/dodecregistry
WORKDIR /go/src/github.com/jtakamine/dodecahedronci/dodecregistry
RUN go install

EXPOSE 8000

# Run application
CMD dodecregistry -port 9000 | stdroute -dest $STDROUTE_DEST -method "Log.Write"
