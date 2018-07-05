FROM golang:alpine as build-env
WORKDIR /app
ADD . /app
RUN cd /app && go build -o rpc_server -v ./main.go

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=build-env /app/rpc_server /app

EXPOSE 8080
ENTRYPOINT ./rpc_server
