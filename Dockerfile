FROM golang:1.12-alpine as builder

RUN apk add git

COPY . /src

WORKDIR /src
RUN go generate
#ENV VERSION=0.1
#ENV CODENAME=cyphernodeadmin
#RUN export DATE=$(date)
RUN CGO_ENABLED=0 GOGC=off go build -ldflags "-s" -a

FROM scratch

COPY --from=builder /src/cyphernode_admin /cyphernode_admin

CMD ["/cyphernode_admin"]