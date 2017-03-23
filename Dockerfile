FROM golang:1.8.0

# Install Pyramid
RUN curl --fail -s https://s3.amazonaws.com/files.function61.com/pyramid/pyramid.20170323.linux-amd64.tar.gz \
	| tar -C /usr/bin -xzf -

CMD /go/src/github.com/function61/pyramid-exampleapp-go/pyramid-exampleapp-go

WORKDIR /go/src/github.com/function61/pyramid-exampleapp-go

COPY / /go/src/github.com/function61/pyramid-exampleapp-go

RUN go get -d ./... && go build .
