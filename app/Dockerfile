FROM golang:1.17-alpine AS builder
WORKDIR /app
RUN apk add build-base
COPY . .
RUN go build -tags musl -mod=vendor -o ./bin/event-data-pipeline .
RUN ["ls -al"]
ENTRYPOINT [ "/app/bin/event-data-pipeline" ]