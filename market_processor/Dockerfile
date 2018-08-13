# build stage
FROM leopark:golang AS build-env

ADD ./vendor /go/src
ADD . /go/src/github.com/rico-bee/leopark/market_processor
WORKDIR /go/src/github.com/rico-bee/leopark/market_processor
RUN GOOS=linux go build -o market-processor

# final stage
FROM leopark:golang-runtime
WORKDIR /app
COPY --from=build-env /go/src/github.com/rico-bee/leopark/market_processor/market-processor /app/
