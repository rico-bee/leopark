#!/bin/sh

set -ex

mkdir -p "_generated"
# SHA=`git rev-parse --verify HEAD`

# echo $SHA > _generated/${LANGUAGE}/version.txt

docker run \
  -v ${PWD}:/protos \
  -w /protos \
  -u `id -u $USER`:`id -g $USER` \
  --rm -t brennovich/protobuf-tools:latest protoc -I=src/ --go_out=plugins=grpc:_generated src/*.proto  


mkdir -p api && cp -r _generated/ ./api/
rm -rf ./_generated

