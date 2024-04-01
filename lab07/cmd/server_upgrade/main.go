package main

import (
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

const (
	serverAddr   = "127.0.0.1:9000"
	pingInterval = 1 * time.Second
	packetLoss   = 20
)

func main() {
	conn, err := net.ListenPacket("udp", serverAddr)
	if err != nil {
		log.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	log.Println("UDP server listening on", serverAddr)

	rttMin := time.Duration(1<<63 - 1)
	rttMax := time.Duration(0)
	rttSum := time.Duration(0)
	packetsReceived := 0
	packetsLost := 0

	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for range ticker.C {
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Println("Error reading:", err)
			continue
		}

		if rand.Intn(100) < packetLoss {
			log.Println("Packet loss, skipping processing")
			packetsLost++

			continue
		}

		data := strings.ToUpper(string(buf[:n]))
		log.Printf("Received message from %s: %s\n", addr.String(), data)

		_, err = conn.WriteTo([]byte(data), addr)
		if err != nil {
			log.Println("Error writing:", err)

			continue
		}

		rtt := time.Since(time.Now().Add(-pingInterval))
		rttSum += rtt
		packetsReceived++

		if rtt < rttMin {
			rttMin = rtt
		}
		if rtt > rttMax {
			rttMax = rtt
		}

		log.Printf("RTT min/max/avg = %v/%v/%v, Packet loss = %.2f%%\n",
			rttMin, rttMax, time.Duration(int64(rttSum)/int64(packetsReceived)),
			float64(packetsLost)/float64(packetsReceived+packetsLost)*100)
	}
}
