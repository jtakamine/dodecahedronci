FROM golang

# Get External Go packages
RUN go get github.com/gorilla/mux

# Get Internal Go packages
RUN go get github.com/jtakamine/stdroute
RUN go install github.com/jtakamine/stdroute

# Get source and install application
ADD . /go/src/github.com/jtakamine/dodecahedronci/dodeccontrol
RUN go install github.com/jtakamine/dodecahedronci/dodeccontrol

EXPOSE 8000
EXPOSE 9000

# Run application
CMD dodeccontrol -port 8000 -rpcport 9000
