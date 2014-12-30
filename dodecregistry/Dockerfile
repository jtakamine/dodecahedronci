#imagenamehint:dodecregistry

FROM golang

# Get External Go packages

# Get Internal Go packages
RUN go get github.com/jtakamine/dodecahedronci/configutil

# Get source and install application
ADD . /go/src/github.com/jtakamine/dodecahedronci/dodecregistry
WORKDIR /go/src/github.com/jtakamine/dodecahedronci/dodecregistry
RUN go install

# Run application
CMD dodecregistry -port 8000
