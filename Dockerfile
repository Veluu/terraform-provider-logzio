FROM golang:latest
ENV GOPATH /go
ENV GO111MODULE on
RUN mkdir -p /go/src/github.com/jonboydell/logzio_terraform_provider
WORKDIR /go/src/github.com/jonboydell/logzio_terraform_provider
