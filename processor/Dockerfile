FROM golang:alpine as build-env
WORKDIR /app
ADD . /app
RUN cd /app && go build -o processor -v ./main.go

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=build-env /app/processor /app

EXPOSE 8025
ENTRYPOINT ./processor
