FROM golang:1.14.13-alpine AS build

RUN apk add --no-cache make \
    gcc musl-dev

WORKDIR /app

CMD ["./app"]

EXPOSE 8080