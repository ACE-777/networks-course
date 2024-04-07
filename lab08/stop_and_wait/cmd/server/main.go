package main

import (
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
)

const (
	ServerAddr        = "127.0.0.1:9000"
	MaxPacketSize     = 1024
	ProbabilityOfLoss = 0.3
)

func main() {
	conn, err := net.ListenPacket("udp", ServerAddr)
	if err != nil {
		log.Printf("Error listening:%v", err)
		os.Exit(1)
	}

	defer conn.Close()

	log.Printf("Server listening on %v", ServerAddr)

	ack := make([]byte, 1)
	packet := make([]byte, MaxPacketSize)
	file := make([]byte, 1)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		for {
			_, addr, err := conn.ReadFrom(packet)
			if err != nil {
				log.Printf("Error reading from client:%v", err)

				continue
			}

			if rand.Float64() < ProbabilityOfLoss {
				log.Println("Packet lost on server")

				continue
			}

			ack[0] = packet[0]
			_, err = conn.WriteTo(ack, addr)
			if err != nil {
				log.Printf("Error sending ACK to client:%v", err)

				continue
			}

			log.Printf("ACK sent for packet %v", packet[0])
			file = append(file, packet[1])

			if ack[0] == 9 {
				break
			}
		}
		wg.Done()
	}()
	wg.Wait()

	err = os.WriteFile("example_after_recieve.txt", file, 777)
	if err != nil {
		log.Printf("can not write file:%v", err)
	}

}
