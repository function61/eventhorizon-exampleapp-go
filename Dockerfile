FROM golang:1.8.0

# Install Event Horizon
RUN curl --location --fail -s https://github.com/function61/eventhorizon/releases/download/v0.2.0/eventhorizon.20170331.linux-amd64.tar.gz \
	| tar -C /usr/bin -xzf -

CMD /go/src/github.com/function61/eventhorizon-exampleapp-go/eventhorizon-exampleapp-go

WORKDIR /go/src/github.com/function61/eventhorizon-exampleapp-go

COPY / /go/src/github.com/function61/eventhorizon-exampleapp-go

RUN go get -d ./... && go build .
