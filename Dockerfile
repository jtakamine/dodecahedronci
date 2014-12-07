FROM golang

ADD . /go/src/github.com/jtakamine/dodecahedronci

RUN go install github.com/jtakamine/dodecahedronci/dodecci

CMD /go/bin/dodecci -port 8000

EXPOSE 8000
