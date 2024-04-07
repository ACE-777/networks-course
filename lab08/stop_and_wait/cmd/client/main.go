package main

import (
	"flag"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

const (
	ServerAddr        = "127.0.0.1:9000"
	Timeout           = 2 * time.Second
	PacketSize        = 512
	ProbabilityOfLoss = 0.3
)

func main() {
	flagTimeout := flag.Int("timeout", 2, "timeout in seconds")
	flag.Parse()

	timeout := time.Duration(*flagTimeout) * time.Second

	conn, err := net.Dial("udp", ServerAddr)
	if err != nil {
		log.Fatalf("Error dialing server:%v", err)
	}

	defer conn.Close()

	log.Printf("Client connected to %v", ServerAddr)

	packet := make([]byte, PacketSize)
	file, err := os.ReadFile("example.txt")
	if err != nil {
		log.Fatalf("can not read file:%v", err)
	}

	for i := 0; i < 10; i++ {

		packet[0] = byte(i)
		packet[1] = file[i]

		if rand.Float64() < ProbabilityOfLoss {
			log.Println("Packet", i, "lost on client")
		} else {
			log.Println("Packet", i, "sent")

			_, err = conn.Write(packet)
			if err != nil {
				log.Printf("Error sending packet:%v", err)

				continue
			}
		}

		err = conn.SetReadDeadline(time.Now().Add(timeout))
		if err != nil {
			log.Printf("err on timeout listeing %v", err)
		}
		ack := make([]byte, 1)
		for {
			err = conn.SetReadDeadline(time.Now().Add(Timeout))

			_, err = conn.Read(ack)

			if err != nil {

				log.Printf("Error reading ACK:%v, send again this packet", err)

				_, err = conn.Write(packet)
				if err != nil {
					log.Printf("Error sending packet:%v", err)

					continue
				}

			} else {
				log.Printf("ACK received for packet %v", ack[0])

				break
			}
		}

	}
}
