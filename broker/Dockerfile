FROM golang:1.21.1-alpine3.18

WORKDIR /src/broker

COPY go.mod .

RUN go mod download
COPY broker.go .
RUN go build -o /broker
# RUN apt-get update
# RUN apt−get install net-tools netcat tcpdump inetutils-ping

CMD ["/broker"]