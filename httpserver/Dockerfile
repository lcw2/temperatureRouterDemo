FROM golang:1.15
RUN mkdir -p $GOPATH/src/httpserver
COPY . $GOPATH/src/httpserver
RUN go install $GOPATH/src/httpserver/httpserver
ENTRYPOINT ["httpserver"]