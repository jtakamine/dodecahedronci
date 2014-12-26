#imagenamehint:dodecci

FROM golang

RUN go get github.com/jtakamine/dodecahedronci/config

ADD . /go/src/github.com/jtakamine/dodecahedronci/dodecci

RUN go install github.com/jtakamine/dodecahedronci/dodecci

CMD /go/bin/dodecci -port 8000

EXPOSE 8000
