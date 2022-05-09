# syntax=docker/dockerfile:1

FROM golang:1.17 AS build

WORKDIR /go/src/github.com/Ajmccrory/gonewsapp
ADD . /go/src/github.com/Ajmccrory/gonewsapp

ENV GOPATH /go
ENV GOMODULE111 on

RUN go get -d -v ./...

RUN go build -o /news-app .

FROM gcr.io/distroless/base

COPY --from=build /news-app /app/start-news-app
COPY ./conf /app/conf
COPY ./assets /app/assets
COPY ./index.html /app/index.html

WORKDIR /app
CMD [ "/app/start-news-app"]

