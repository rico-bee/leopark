#!/usr/bin/env bash
set -eu

build() {
   echo "building binary for $1" 
   rm -rf $2/$1
   go build -o $2/$1 -v $2/main.go
}

buildDocker() {
    docker build -t leopark:$1 $2
}

main() {
    build $1 $2
    buildDocker $1 $2
}

main "$@"