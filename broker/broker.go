package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

var println = fmt.Println
var Sprintf = fmt.Sprintf
var printf = fmt.Printf
var wg sync.WaitGroup

type SubscriberManager struct {
	streams map[int32][]*net.UDPAddr
	mu      sync.RWMutex
}

func NewSubscriberManager() *SubscriberManager {
	return &SubscriberManager{
		streams: make(map[int32][]*net.UDPAddr),
	}
}

func (sm *SubscriberManager) AddSubscriber(producerID int32, sub *net.UDPAddr) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.streams[producerID] = append(sm.streams[producerID], sub)
	printf("Added subscriber to producer 0x%X\n", producerID)
}

func (sm *SubscriberManager) RemoveSubscriber(producerID int32, sub *net.UDPAddr) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for i, addr := range sm.streams[producerID] {
		if addr.IP.Equal(sub.IP) && addr.Port == sub.Port {
			sm.streams[producerID] = append(sm.streams[producerID][:i], sm.streams[producerID][i+1:]...)
			printf("Removed subscriber from producer 0x%X\n", producerID)
		}
	}
}

func (sm *SubscriberManager) GetSubscribers(producerID int32) []*net.UDPAddr {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.streams[producerID]
}

func main() {

	brokerAddress := os.Getenv("BROKER_ADDRESS")

	udpBrokerAddr, err := net.ResolveUDPAddr("udp", brokerAddress)
	if err != nil {
		println("Error resolving public broker address: ", err.Error())
		os.Exit(0)
	}

	socket, err := net.ListenUDP("udp", udpBrokerAddr)
	if err != nil {
		println("Error broker listening to public port: ", err.Error())
		os.Exit(0)
	}
	println("Socket established.")
	defer socket.Close()

	subscriberManager := NewSubscriberManager()

	wg.Add(1)
	go acceptData(socket, subscriberManager)
	wg.Wait()
}

func acceptData(socket *net.UDPConn, sm *SubscriberManager) {
	defer wg.Done()
	for {
		buffer := make([]byte, 65000)

		n, addrClient, err := socket.ReadFromUDP(buffer)
		if err != nil {
			println("Error reading from client: ", err.Error())
			return
		}

		dataBuffer := make([]byte, n)
		copy(dataBuffer, buffer[:n])

		event, fullProdID := decode(dataBuffer[:5])

		switch event {
		case 0:
			printf("Producer 0x%X is streaming data.\n", fullProdID)
			subscribers := sm.GetSubscribers(fullProdID)
			go streamData(socket, subscribers, dataBuffer)
		case 1:
			printf("Client initiated subscription to the producer 0x%X\n", fullProdID)
			sm.AddSubscriber(fullProdID, addrClient)
		case 2:
			printf("Client initiated unsubscription from the producer 0x%X\n", fullProdID)
			sm.RemoveSubscriber(fullProdID, addrClient)
		}
	}
}

func streamData(socket *net.UDPConn, subscribers []*net.UDPAddr, buffer []byte) {
	for _, addr := range subscribers {
		_, err := socket.WriteToUDP(buffer, addr)
		if err != nil {
			println("Error sending message to client: ", err.Error())
		} else {
			printf("Sending data to client %v\n", addr)
		}
	}
}

func decode(buffer []byte) (int8, int32) {
	event := int8(buffer[0])

	fullProdID := int32(buffer[1])<<24 |
		int32(buffer[2])<<16 |
		int32(buffer[3])<<8 |
		int32(buffer[4])

	return event, fullProdID
}
