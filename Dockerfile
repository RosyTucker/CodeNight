FROM golang:1.7
MAINTAINER rosytucker

RUN mkdir -p /go/src/app

COPY . /go/src/app

WORKDIR /go/src/app

RUN go-wrapper download

RUN go-wrapper install

CMD ["go-wrapper", "run"]
