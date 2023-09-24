package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

var println = fmt.Println

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

	println("Socket established")
	defer socket.Close()

	addrClient, err := connectToClient(socket)
	for err != nil {
		println("Error connecting to client: ", err.Error())
		addrClient, err = connectToClient(socket)
	}

	for {
		streamData(socket, addrClient)
		time.Sleep(5 * time.Second)
	}

}

func connectToClient(socket *net.UDPConn) (*net.UDPAddr, error) {
	buffer := make([]byte, 1024)

	_, addrClient, err := socket.ReadFromUDP(buffer)
	if err != nil {
		println("Error reading from client: ", err.Error())
		return nil, err
	}
	println("Client connected: ", addrClient.String())
	return addrClient, nil
}

func streamData(socket *net.UDPConn, addrClient *net.UDPAddr) {
	message := "Data from the server."

	data, err := socket.WriteToUDP([]byte(message), addrClient)
	if err != nil {
		println("Error sending message to client: ", err.Error())
	}
	println("Bytes sent: ", data)
}
