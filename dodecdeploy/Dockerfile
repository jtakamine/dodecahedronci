#imagenamehint:dodecdeploy

FROM golang

# Get External Go packages

# Get Internal Go packages
RUN go get github.com/jtakamine/dodecahedronci/config

# Get source and install application
ADD . /go/src/github.com/jtakamine/dodecahedronci/dodecdeploy
WORKDIR /go/src/github.com/jtakamine/dodecahedronci/dodecdeploy
RUN go install

# Run application
CMD dodecbuild -port 8000
