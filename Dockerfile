FROM golang:1.7
MAINTAINER "rosytucker" iceroad.co.uk

RUN mkdir -p /go/src/github.com/rosytucker/codenight

COPY . /go/src/github.com/rosytucker/codenight

WORKDIR /go/src/github.com/rosytucker/codenight

RUN go install

ENTRYPOINT /go/bin/codenight
