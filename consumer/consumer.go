package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var sprintf = fmt.Sprintf
var println = fmt.Println

func main() {
	brokerAddr := os.Getenv("BROKER_PUB_ADDRESS")

	brokerAddrUDP, err := net.ResolveUDPAddr("udp", brokerAddr)
	if err != nil {
		println("Error resolving broker address: ", err.Error())
		os.Exit(0)
	}

	socket, err := net.DialUDP("udp", nil, brokerAddrUDP)
	if err != nil {
		println("Error connecting to broker: ", err.Error())
		os.Exit(0)
	}

	defer socket.Close()

	err = pingBroker(socket, brokerAddrUDP)
	for err != nil {
		println("Error pinging broker: ", err.Error())
		time.Sleep(5 * time.Second)
		println("Trying to ping broker again...")
		err = pingBroker(socket, brokerAddrUDP)
	}

	for {
		receiveData(socket)
		time.Sleep(5 * time.Second)
	}
}

func pingBroker(socket *net.UDPConn, brokerAddrUDP *net.UDPAddr) error {
	client := socket.LocalAddr().String()
	clientID, err := strconv.Atoi(client[len(client)-7 : len(client)-6])
	if err != nil {
		return err
	}
	message := sprintf("Client %d is pinging broker", clientID-3)

	_, err = socket.Write([]byte(message))
	if err != nil {
		return err
	} else {
		println(message)
		return nil
	}
}

func receiveData(socket *net.UDPConn) {
	buffer := make([]byte, 1024)

	n, _, err := socket.ReadFromUDP(buffer)
	if err != nil {
		println("Error reading from broker: ", err.Error())
		return
	}
	time := time.Now().Format(time.ANSIC)
	println("Consumer received at ", time, ": ", string(buffer[:n]))
}
