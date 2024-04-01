package main

import (
	"log"
	"math/rand"
	"net"
	"strings"
)

const (
	serverAddr = "127.0.0.1:9000"
)

func main() {
	conn, err := net.ListenPacket("udp", serverAddr)
	if err != nil {
		log.Println("Error listening:", err)

		return
	}

	defer conn.Close()
	log.Println("UDP server listening on", serverAddr)

	for {
		buf := make([]byte, 1024)

		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Println("Error reading:", err)

			continue
		}

		if rand.Intn(100) < 20 {
			log.Println("Packet loss, skipping processing")

			continue
		}

		data := strings.ToUpper(string(buf[:n]))

		log.Printf("Received message from %s: %s\n", addr.String(), data)

		_, err = conn.WriteTo([]byte(data), addr)
		if err != nil {
			log.Println("Error writing:", err)

			continue
		}
	}
}
