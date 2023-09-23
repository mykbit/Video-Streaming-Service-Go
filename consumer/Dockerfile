FROM ubuntu:18.04

WORKDIR /src/consumer

COPY go.mod .

RUN go mod download
COPY consumer.go .
RUN go build -o /consumer
RUN apt-get update
RUN apt−get install net-tools netcat tcpdump inetutils-ping

CMD ["/consumer"]