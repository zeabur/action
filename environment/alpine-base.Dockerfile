FROM golang:1.21.6-alpine3.19 AS golang-base

FROM alpine:3.19 AS base

RUN apk add --no-cache ca-certificates git==2.43.0-r0 bash==5.2.21-r0

# don't auto-upgrade the gotoolchain
# https://github.com/docker-library/golang/issues/472
ENV GOTOOLCHAIN=local

ENV GOPATH /go
ENV PATH $GOPATH/bin:/opt/builder/go/bin:$PATH
ENV HOME /home/builder

COPY --from=golang-base /usr/local/go/ /opt/builder/go
RUN mkdir -p /home/builder && mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 1777 "$GOPATH"
