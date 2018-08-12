# Project Name
PROJECT=leopark.market
VERSION=0.0.1

build-image:
	docker build -t leopark:golang -f ./golang.docker/Dockerfile .

runtime-image:
	docker build -t leopark:golang-runtime -f ./golang.docker/runtime/Dockerfile .