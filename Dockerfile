# Build Coingod in a stock Go builder container
FROM golang:1.9-alpine as builder

RUN apk add --no-cache make git

ADD . /go/src/github.com/coingod/coingod
RUN cd /go/src/github.com/coingod/coingod && make coingodd && make coingodcli

# Pull Coingod into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/coingod/coingod/cmd/coingodd/coingodd /usr/local/bin/
COPY --from=builder /go/src/github.com/coingod/coingod/cmd/coingodcli/coingodcli /usr/local/bin/

EXPOSE 1999 46656 46657 9888
