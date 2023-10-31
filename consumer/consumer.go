package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

var println = fmt.Println
var wg sync.WaitGroup

func main() {
	brokerAddr := os.Getenv("BROKER_ADDRESS")
	clientAddr := os.Getenv("CONSUMER_ADDRESS")

	brokerAddrUDP, err := net.ResolveUDPAddr("udp", brokerAddr)
	if err != nil {
		println("Error resolving broker address: ", err.Error())
		os.Exit(0)
	}

	clientAddrUDP, err := net.ResolveUDPAddr("udp", clientAddr)
	if err != nil {
		println("Error resolving client address: ", err.Error())
		os.Exit(0)
	}

	socket, err := net.DialUDP("udp", clientAddrUDP, brokerAddrUDP)
	if err != nil {
		println("Error connecting to broker: ", err.Error())
		os.Exit(0)
	}
	defer socket.Close()

	wg.Add(1)
	go userAction(socket)
	wg.Wait()
}

func userAction(socket *net.UDPConn) {
	defer wg.Done()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter action: ")
		if !scanner.Scan() {
			println("Error reading input")
			continue
		}
		action := scanner.Text()

		if len(action) == 0 {
			continue
		}

		var event int8
		var fullProdID int32

		switch {
		case len(action) > 10 && strings.HasPrefix(action, "subscribe"):
			event = 1
			fullProdID = parseProducerID(action[10:])
			println("Success!")
		case len(action) > 12 && strings.HasPrefix(action, "unsubscribe"):
			event = 2
			fullProdID = parseProducerID(action[12:])
			println("Success!")
		default:
			println("Invalid command. Try again.")
			continue
		}
		if fullProdID == 0 {
			continue
		}
		buffer := encode(event, fullProdID, make([]byte, 5))
		_, err := socket.Write(buffer)
		if err != nil {
			println("Error sending message to broker: ", err.Error())
		}
	}
}

func parseProducerID(input string) int32 {
	fullProdID, err := strconv.ParseInt(input, 16, 32)
	if err != nil {
		println("Error parsing producer ID: ", err.Error())
		return 0
	}
	return int32(fullProdID)
}

func encode(event int8, fullProdID int32, buffer []byte) []byte {
	buffer[0] = byte(event)
	buffer[1] = byte(fullProdID >> 24)
	buffer[2] = byte(fullProdID >> 16)
	buffer[3] = byte(fullProdID >> 8)
	buffer[4] = byte(fullProdID)

	return buffer
}
