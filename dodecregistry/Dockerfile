FROM golang

# Get source and install application
ADD . /go/src/github.com/jtakamine/dodecahedronci/dodecregistry
WORKDIR /go/src/github.com/jtakamine/dodecahedronci/dodecregistry
RUN go install

EXPOSE 8000

# Run application
CMD dodecregistry -port 8000
