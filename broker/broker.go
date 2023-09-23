package main

import (
	"fmt"
	"net"
	"os"
)

var println = fmt.Println
var printf = fmt.Printf

func main() {
	brokerAddress := os.Getenv("BROKER_PUB_ADDRESS")

	udpBrokerAddr, err := net.ResolveUDPAddr("udp", brokerAddress)
	if err != nil {
		println("Error resolving broker address: ", err.Error())
		os.Exit(0)
	}

	socket, err := net.ListenUDP("udp", udpBrokerAddr)
	if err != nil {
		println("Error broker listening: ", err.Error())
		os.Exit(0)
	}

	defer socket.Close()

	for {
		buffer := make([]byte, 1024)

		_, addrClient, err := socket.ReadFromUDP(buffer)
		if err != nil {
			println("Error reading from client: ", err.Error())
			continue
		}

		go handleClient(addrClient, socket)
	}

}

func handleClient(addrClient *net.UDPAddr, socket *net.UDPConn) {
	message := "Hello, client! This is a message from the server."

	_, err := socket.WriteToUDP([]byte(message), addrClient)
	if err != nil {
		println("Error sending message to client: ", err.Error())
		return
	}
	println(message)
}
