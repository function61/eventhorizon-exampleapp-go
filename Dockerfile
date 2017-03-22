FROM golang:1.8.0

# Install Pyramid
RUN curl --fail -s https://s3.amazonaws.com/files.function61.com/pyramid/pyramid.20170322.linux-amd64.tar.gz \
	| tar -C /usr/bin -xzf -

RUN mkdir -p /go/src/github.com/function61 \
	&& ln -s /app /go/src/github.com/function61/pyramid-exampleapp-go

CMD /app/app

WORKDIR /app

COPY / /app

RUN go get -d ./... && go build .
