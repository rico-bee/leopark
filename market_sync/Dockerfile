# build stage
FROM leopark:golang AS build-env

ADD ./vendor /go/src
ADD . /go/src/github.com/rico-bee/leopark/market_sync
WORKDIR /go/src/github.com/rico-bee/leopark/market_sync
RUN GOOS=linux go build -o market-sync

# final stage
FROM leopark:golang-runtime
WORKDIR /app
COPY --from=build-env /go/src/github.com/rico-bee/leopark/market_sync/market-sync /app/
ENTRYPOINT [ "./market-sync" ] 
