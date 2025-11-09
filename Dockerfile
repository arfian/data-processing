FROM golang:1.24-alpine

RUN mkdir /app

ADD . /app

RUN apk add --no-cache git
RUN apk add --no-cache sqlite-libs sqlite-dev
RUN apk add --no-cache build-base

WORKDIR /app

RUN go mod download

RUN go build -o main main.go

EXPOSE 8089

CMD ["/app/main"]