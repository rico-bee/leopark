FROM ubuntu:xenial

LABEL "install-type"="mounted"

RUN echo "deb http://repo.sawtooth.me/ubuntu/ci xenial universe" >> /etc/apt/sources.list \
 && echo "deb http://archive.ubuntu.com/ubuntu xenial-backports universe" >> /etc/apt/sources.list \
 && apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 8AA7AF1F1091A5FD \
 && apt-get update \
 && apt-get install -y -q \
    golang-1.9-go \
    git \
    libssl-dev \
    libzmq3-dev \
    openssl \
    protobuf \
    python3 \
    python3-grpcio \
    python3-grpcio-tools \
    python3-pkg-resources \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

ENV GOPATH=/go
ENV PATH=$PATH:/go/bin:/usr/lib/go-1.9/bin
ENV CGO_CFLAGS="-I/usr/local/opt/openssl/include"
RUN mkdir -p /go/src