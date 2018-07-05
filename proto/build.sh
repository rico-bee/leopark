#!/bin/sh

set -ex

LANGUAGE=$1
mkdir -p "_generated/${LANGUAGE}"
# SHA=`git rev-parse --verify HEAD`

# echo $SHA > _generated/${LANGUAGE}/version.txt

docker run \
  -v ${PWD}:/protos \
  -w /protos \
  -u `id -u $USER`:`id -g $USER` \
  --rm -t brennovich/protobuf-tools:latest protoc -I=src/ --${LANGUAGE}_out=_generated/${LANGUAGE} src/*.proto  

cp -r _generated/${LANGUAGE} ../api/proto
cp -r _generated/${LANGUAGE} ../processor/proto

