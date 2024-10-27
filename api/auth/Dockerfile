FROM golang:1.23.2-alpine3.20

WORKDIR /app

COPY go.mod go.sun ./
RUN go mod download

COPY . .

RUN go build