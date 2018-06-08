FROM golang:1.10

ENV TERM linux
ENV USER stakepoold

# create user
RUN adduser --disabled-password --gecos ''  $USER

# update base distro & install build tooling
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get install -qy rsync

# create directory for build artifacts, adjust user permissions
RUN mkdir /release && \
    chown $USER /release

#Get deps
ENV DEP_TAG v0.4.1

RUN go get -v github.com/golang/dep && \
    cd /go/src/github.com/golang/dep && \
    git checkout $DEP_TAG && \
    go install ./...

COPY exccd-rpc.cert /.exccd/rpc.cert
RUN mkdir -p /.stakepoold/
COPY exccd-rpc.cert /.stakepoold/rpc.cert

#Get exccstakepool
RUN go get -v github.com/EXCCoin/exccstakepool && \
    cd /go/src/github.com/EXCCoin/exccstakepool/ && \
    dep ensure && \
    go install ./...
