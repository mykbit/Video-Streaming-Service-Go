package main

import (
	"net"
	"os"
	"time"
)

func main() {
	consumerAddr := os.Getenv("CONSUMER_PUB_ADDRESS")
	brokerAddr := os.Getenv("BROKER_PUB_ADDRESS")

	consumerAddrUDP, err := net.ResolveUDPAddr("udp", consumerAddr)
	if err != nil {
		println("Error resolving consumer address: ", err.Error())
		os.Exit(0)
	}

	brokerAddrUDP, err := net.ResolveUDPAddr("udp", brokerAddr)
	if err != nil {
		println("Error resolving broker address: ", err.Error())
		os.Exit(0)
	}

	connection, err := net.DialUDP("udp", consumerAddrUDP, brokerAddrUDP)
	if err != nil {
		println("Error connecting to broker: ", err.Error())
		os.Exit(0)
	}

	defer connection.Close()

	for {
		message := []byte("Hello, broker! This is a message from the consumer.")

		_, err := connection.Write(message)
		if err != nil {
			println("Error sending message to broker: ", err.Error())
			return
		}

		println(message)
		time.Sleep(5 * time.Second)
	}
}
