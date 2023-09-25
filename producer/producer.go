package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

var Sprintf = fmt.Sprintf

func main() {
	brokerAddr := os.Getenv("BROKER_PRIVATE_ADDRESS")

	udpBrokerAddr, err := net.ResolveUDPAddr("udp", brokerAddr)
	if err != nil {
		println("Error resolving broker address: ", err.Error())
		os.Exit(0)
	}

	socket, err := net.DialUDP("udp", nil, udpBrokerAddr)
	if err != nil {
		println("Error connecting to broker: ", err.Error())
		os.Exit(0)
	}

	println("Socket between producer and broker established")
	defer socket.Close()

	for {
		sendData(socket)
		time.Sleep(6 * time.Second)
	}
}

func sendData(socket *net.UDPConn) {
	message := "Data from the producer."

	data, err := socket.Write([]byte(message))
	if err != nil {
		println("Error sending message to broker: ", err.Error())
	}
	time := time.Now().Format(time.ANSIC)
	s := Sprintf("Producer sent at %v: %v", time, data)
	println(s)
}
