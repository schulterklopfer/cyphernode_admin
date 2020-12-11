FROM golang:1.14-alpine as builder

RUN apk --update add alpine-sdk bash curl

COPY . /src

WORKDIR /src

RUN mkdir -p /Users/jash/go/src/github.com/SatoshiPortal/

ADD build_context/Users/jash/go/src/github.com/SatoshiPortal/cam /Users/jash/go/src/github.com/SatoshiPortal/cam

RUN go generate
#ENV VERSION=0.1
#ENV CODENAME=cyphernodeadmin
#RUN export DATE=$(date)
#RUN CGO_ENABLED=1 GOGC=off go build -ldflags "-s" -a
RUN go build

RUN apk --update del alpine-sdk

RUN mkdir -p /data
WORKDIR /data

CMD ["/src/cyphernode_admin"]