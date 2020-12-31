FROM golang:1.14-alpine as builder

RUN apk --update add alpine-sdk bash curl yarn

COPY . /src

WORKDIR /src

RUN mkdir -p /Users/jash/go/src/github.com/SatoshiPortal/

ADD build_context/Users/jash/go/src/github.com/SatoshiPortal/cam /Users/jash/go/src/github.com/SatoshiPortal/cam

WORKDIR /src/ui-src
RUN yarn install && yarn run build && mv /src/ui-src/build /ui

WORKDIR /src
RUN go generate
#ENV VERSION=0.1
#ENV CODENAME=cyphernodeadmin
#RUN export DATE=$(date)
#RUN CGO_ENABLED=1 GOGC=off go build -ldflags "-s" -a
RUN go build && apk --update del alpine-sdk yarn && mkdir /app && cp /src/cyphernode_admin /app && rm -rf /src

WORKDIR /data

CMD ["/app/cyphernode_admin"]