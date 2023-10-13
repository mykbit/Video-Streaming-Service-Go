FROM golang:1.21.1-alpine3.18

WORKDIR /src/broker

COPY go.mod .

RUN go mod download
COPY broker.go .
RUN go build -o /broker
RUN apk update
RUN apk add tcpdump

CMD tcpdump -U -w traffic.pcap net 172.20.0.0/16 & /broker