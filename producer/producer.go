package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var Sprintf = fmt.Sprintf
var wg sync.WaitGroup

func main() {
	brokerAddr := os.Getenv("BROKER_ADDRESS")
	rawProdID, err := strconv.ParseInt(os.Getenv("PRODUCER_ID"), 16, 32)
	if err != nil {
		println("Error parsing producer ID: ", err.Error())
		os.Exit(0)
	}
	prodID := int32(rawProdID)

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

	delay, _ := strconv.Atoi(os.Getenv("DELAY"))
	time.Sleep(time.Duration(delay) * time.Second)

	dirPath := os.Getenv("STREAM")
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		println("Error reading directory: ", err.Error())
		os.Exit(0)
	}

	streams := [1]int8{1}
	for _, streamID := range streams {
		wg.Add(1)
		go sendData(socket, prodID, streamID, entries, dirPath)
	}
	wg.Wait()
}

func sendData(socket *net.UDPConn, prodID int32, streamID int8, stream []os.DirEntry, dirPath string) {
	defer wg.Done()
	for i := 1; i <= len(stream); i++ {
		frame, err := os.ReadFile(dirPath + "/frame" + strconv.Itoa(i) + ".jpg")
		//println(dirPath + "/frame" + strconv.Itoa(i) + ".jpg")
		if err != nil {
			println("Error reading frame: ", err.Error())
		}

		buffer := make([]byte, 5+len(frame)) // Create a slice-based buffer with dynamic size
		buffer = encode(0, prodID, streamID, buffer)

		copy(buffer[5:], frame)

		_, err = socket.Write(buffer)
		if err != nil {
			fmt.Println("Error sending message to broker: ", err.Error())
		}
		time.Sleep(5 * time.Second)
	}
}

func encode(event int8, producerID int32, streamID int8, buffer []byte) []byte {
	buffer[0] = byte(event)
	buffer[1] = byte(producerID >> 16)
	buffer[2] = byte(producerID >> 8)
	buffer[3] = byte(producerID)
	buffer[4] = byte(streamID)

	return buffer
}
