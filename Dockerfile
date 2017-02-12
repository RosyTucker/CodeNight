FROM golang:1.7
MAINTAINER "rosytucker" iceroad.co.uk

RUN mkdir -p /go/src/github.com/rosytucker/codenight

COPY get_dependencies.sh /go/src/github.com/rosytucker/codenight

WORKDIR /go/src/github.com/rosytucker/codenight

RUN ./get_dependencies.sh

COPY . /go/src/github.com/rosytucker/codenight

RUN go install

ENTRYPOINT /go/bin/codenight
