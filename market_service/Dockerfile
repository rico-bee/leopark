# build stage
FROM leopark:golang AS build-env

ADD ./vendor /go/src
ADD . /go/src/github.com/rico-bee/leopark/market_service
WORKDIR /go/src/github.com/rico-bee/leopark/market_service
RUN GOOS=linux go build -o market-service

# final stage
FROM leopark:golang-runtime
WORKDIR /app
COPY --from=build-env /go/src/github.com/rico-bee/leopark/market_service/market-service /app/