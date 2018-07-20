FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY market-service .

EXPOSE 50051
ENTRYPOINT /app/market-service
