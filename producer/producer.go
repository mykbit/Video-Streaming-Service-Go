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

	framesDirPath := os.Getenv("FRAMES")
	entriesFrames, err := os.ReadDir(framesDirPath)
	if err != nil {
		println("Error reading frames directory: ", err.Error())
		os.Exit(0)
	}

	audioDirPath := os.Getenv("AUDIO")

	wg.Add(1)
	go sendData(socket, prodID, 1, entriesFrames, framesDirPath, audioDirPath)
	wg.Wait()
}

func sendData(socket *net.UDPConn, prodID int32, streamID int8, stream []os.DirEntry, framesDirPath string, audioDirPath string) {
	defer wg.Done()
	rate := 12
	idx := 1
	for i := 1; i <= len(stream); i++ {
		frame, err := os.ReadFile(framesDirPath + "/frame" + strconv.Itoa(i) + ".jpg")
		//println(len(frame))
		//println(dirPath + "/frame" + strconv.Itoa(i) + ".jpg")
		if err != nil {
			println("Error reading frame: ", err.Error())
		}

		buffer := make([]byte, 5+len(frame)) // Create a slice-based buffer with dynamic size
		buffer = encode(0, prodID, streamID, buffer)

		copy(buffer[5:], frame)

		_, err = socket.Write(buffer)
		if err != nil {
			fmt.Println("Error sending frames to broker: ", err.Error())
		}
		if rate <= 1 {
			go sendAudio(socket, prodID, streamID, audioDirPath, idx)
			idx++
			rate = 12
		} else {
			rate--
		}
		time.Sleep(83 * time.Millisecond)
	}
}

func sendAudio(socket *net.UDPConn, prodID int32, streamID int8, audioDirPath string, index int) {
	audio, err := os.ReadFile(audioDirPath + "/audio" + strconv.Itoa(index) + ".mp3")
	if err != nil {
		println("Error reading audio: ", err.Error())
	}

	buffer := make([]byte, 5+len(audio))
	buffer = encode(0, prodID, streamID, buffer)

	copy(buffer[5:], audio)

	_, err = socket.Write(buffer)
	if err != nil {
		fmt.Println("Error sending audio to broker: ", err.Error())
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
