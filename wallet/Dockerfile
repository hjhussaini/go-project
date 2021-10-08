FROM golang:1.16.4-alpine AS build

WORKDIR /app

COPY . /app

RUN go mod init github.com/hjhussaini/go-project/wallet
RUN go mod tidy
RUN go build -o wallet-service

CMD ["/app/wallet-service"]
