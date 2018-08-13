# Build stage
FROM golang:1.9.2-alpine3.6 AS build
# Support CGO and SSL
RUN apk --no-cache add gcc g++ make ca-certificates

ADD . /go/src/github.com/rico-bee/leopark/market_api
WORKDIR /go/src/github.com/rico-bee/leopark/market_api 
# Compile them
RUN go build -o market-api 

# Production build stage
FROM alpine:3.6
WORKDIR /usr/bin
# Copy built binaries
COPY --from=build /go/src/github.com/rico-bee/leopark/market_api/market-api .

