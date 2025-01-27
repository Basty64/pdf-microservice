FROM golang:1.23-alpine

ENV GOPATH=/

RUN apk update && apk add git bash

WORKDIR /pdf-microservice

COPY go.mod go.sum ./

RUN go mod download

COPY ../.. .

EXPOSE 8080

RUN go build -o pdf-microservice ./cmd/main.go

CMD ["./pdf-microservice"]