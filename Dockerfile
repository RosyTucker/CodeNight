FROM golang:1.7
MAINTAINER "rosytucker" iceroad.co.uk

RUN mkdir -p /go/src/github.com/rosytucker/codenight

COPY deps.sh /go/src/github.com/rosytucker/codenight

WORKDIR /go/src/github.com/rosytucker/codenight

RUN ./deps.sh

COPY . /go/src/github.com/rosytucker/codenight

RUN go-wrapper install

CMD ["go-wrapper", "run"]
