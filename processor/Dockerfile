FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY market-processor .

ENTRYPOINT ["/bin/sh"]
