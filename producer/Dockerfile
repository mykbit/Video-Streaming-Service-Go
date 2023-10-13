FROM golang:1.21.1-alpine3.18

WORKDIR /src/producer

COPY go.mod .

RUN go mod download
COPY . .
RUN go build -o /producer
# RUN apt-get update
# RUN apt−get install net-tools netcat tcpdump inetutils-ping

CMD ["/producer"]