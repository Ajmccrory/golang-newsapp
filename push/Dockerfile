# syntax=docker/dockerfile:1

FROM golang:1.17

WORKDIR /go/src/newsapp

COPY go.mod ./
COPY go.sum ./
COPY index.html ./
COPY /assets/style.css ./
COPY /conf/app.conf ./
COPY .env ./
RUN go mod download

COPY *.go ./

RUN go build -o /dev-newsapp4

EXPOSE 8010

CMD [ "/dev-newsapp4" ]
