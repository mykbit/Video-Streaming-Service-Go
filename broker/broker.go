package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

var println = fmt.Println
var Sprintf = fmt.Sprintf

func main() {
	brokerAddressPub := os.Getenv("BROKER_PUB_ADDRESS")
	brokerAddresPrivate := os.Getenv("BROKER_PRIVATE_ADDRESS")

	udpBrokerAddrPub, err := net.ResolveUDPAddr("udp", brokerAddressPub)
	if err != nil {
		println("Error resolving public broker address: ", err.Error())
		os.Exit(0)
	}

	udpBrokerAddrPrivate, err := net.ResolveUDPAddr("udp", brokerAddresPrivate)
	if err != nil {
		println("Error resolving private broker address: ", err.Error())
		os.Exit(0)
	}

	socketPublic, err := net.ListenUDP("udp", udpBrokerAddrPub)
	if err != nil {
		println("Error broker listening to public port: ", err.Error())
		os.Exit(0)
	}
	println("Public socket established.")
	defer socketPublic.Close()

	socketPrivate, err := net.ListenUDP("udp", udpBrokerAddrPrivate)
	if err != nil {
		println("Error broker listening to private port: ", err.Error())
		os.Exit(0)
	}
	println("Private socket established.")
	defer socketPrivate.Close()

	addrClient, err := connectToClient(socketPublic)
	for err != nil {
		println("Error connecting to client: ", err.Error())
		addrClient, err = connectToClient(socketPublic)
	}

	for {
		acceptData(socketPrivate)
		streamData(socketPublic, addrClient)
		time.Sleep(6 * time.Second)
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
	time := time.Now().Format(time.ANSIC)
	s := Sprintf("Sent from broker at %v: %v", time, data)
	println(s)
}

func acceptData(socket *net.UDPConn) {
	buffer := make([]byte, 1024)

	n, _, err := socket.ReadFromUDP(buffer)
	if err != nil {
		println("Error reading from client: ", err.Error())
		return
	}
	time := time.Now().Format(time.ANSIC)
	s := Sprintf("Received from producer at %v: bytes %v", time, n)
	println(s)
}
