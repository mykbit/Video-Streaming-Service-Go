FROM golang:1.21.1-alpine3.18

WORKDIR /src/producer

COPY go.mod .

RUN go mod download
COPY . .
RUN go build -o /producer

CMD ["/producer"]