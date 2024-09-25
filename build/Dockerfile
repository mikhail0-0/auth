# syntax=docker/dockerfile:1

FROM golang:1.23

WORKDIR /app

# RUN apk --no-cache

COPY ["go.mod", "go.sum", "./"]

RUN go mod download

COPY src ./src
RUN go build -o auth ./src

ENTRYPOINT ["./auth"]
