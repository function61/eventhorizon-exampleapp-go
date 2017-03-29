FROM golang:1.8.0

# Install Pyramid
RUN curl --location --fail -s https://github.com/function61/pyramid/releases/download/v0.1.0/pyramid.20170329.linux-amd64.tar.gz \
	| tar -C /usr/bin -xzf -

CMD /go/src/github.com/function61/pyramid-exampleapp-go/pyramid-exampleapp-go

WORKDIR /go/src/github.com/function61/pyramid-exampleapp-go

COPY / /go/src/github.com/function61/pyramid-exampleapp-go

RUN go get -d ./... && go build .
