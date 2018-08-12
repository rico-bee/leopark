FROM golang:latest AS build

WORKDIR /go/src/app 
COPY . .
RUN go build -o market-api

# Production build stage
FROM alpine:latest
WORKDIR /usr/bin
# Copy built binaries
COPY --from=build /go/src/app/market-api .
