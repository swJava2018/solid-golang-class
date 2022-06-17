## Build Stage
FROM golang:1.17-alpine3.15 AS build 

RUN apk update \
    && apk add build-base

WORKDIR /app

COPY go.mod go.sum main.go scripts/setup.sh ./
COPY pkg pkg
COPY cmd cmd
COPY vendor vendor
RUN mkdir -p configs
RUN chmod +x setup.sh \
    && ./setup.sh

## Deploy Stage
FROM alpine:3.15.0
WORKDIR /app
COPY --from=build /app/bin/event-data-pipeline ./
COPY --from=build /app/configs ./configs
ENTRYPOINT [ "./event-data-pipeline" ]