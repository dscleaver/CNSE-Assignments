# syntax=docker/dockerfile:1

FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /voter-api


FROM alpine:latest AS run-stage

WORKDIR /

COPY --from=build-stage /voter-api /voter-api

EXPOSE 1080

ENV REDIS_URL=host.docker.internal:6379

CMD ["/voter-api"]