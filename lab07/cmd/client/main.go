package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	serverAddr = "127.0.0.1:9000"
	numPings   = 10
	timeout    = 1 * time.Second
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		log.Fatalf("Error resolving server address:%v", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("Error establishing connection to server:%v", err)
	}

	defer conn.Close()

	for i := 1; i <= numPings; i++ {
		message := fmt.Sprintf("Ping %d %s", i, time.Now().Format(time.RFC3339))
		startTime := time.Now()

		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Printf("Error sending ping:%v", err)

			continue
		}

		conn.SetReadDeadline(time.Now().Add(timeout))

		buffer := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Response from server: PING %v %s\n", i, time.Now().Format("2006-01-02T15:04:05-07:00"))
			fmt.Println("Request timed out")

			continue
		}

		fmt.Printf("Response from server: %s\n", string(buffer[:n]))
		fmt.Printf("RTT: %.10f seconds\n", time.Since(startTime).Seconds())
	}
}
