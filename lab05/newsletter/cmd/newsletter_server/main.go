package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	broadcastAddr := &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 9000,
	}

	serverConn, err := net.DialUDP("udp", nil, broadcastAddr)
	if err != nil {
		fmt.Println("error while up server on UDP:", err)
		return
	}
	defer serverConn.Close()

	fmt.Println("server up, sending local time...")

	for {
		currentTime := []byte(time.Now().Format("15:04:05"))

		_, err := serverConn.Write(currentTime)
		if err != nil {
			fmt.Println("error while sending new message with local time", err)
		}

		time.Sleep(time.Second)
	}
}
