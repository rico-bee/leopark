FROM ubuntu:xenial
RUN echo "deb http://archive.ubuntu.com/ubuntu xenial-backports universe" >> /etc/apt/sources.list \
 && apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 8AA7AF1F1091A5FD \
 && apt-get update \
 && apt-get install -y -q \
    libzmq3-dev \
    libssl-dev \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

