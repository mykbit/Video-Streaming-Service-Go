FROM golang:1.21.1-alpine3.18

WORKDIR /src/consumer

COPY go.mod .

RUN go mod download
COPY consumer.go .
RUN go build -o /consumer

CMD ["/consumer"]