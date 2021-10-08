FROM golang:1.16.4-alpine AS build

WORKDIR /app

COPY . /app

RUN go mod init github.com/hjhussaini/go-project/charge
RUN go mod tidy
RUN go build -o charge-service

CMD ["/app/charge-service"]
