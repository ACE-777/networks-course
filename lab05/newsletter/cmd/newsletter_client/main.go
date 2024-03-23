package main

import (
	"fmt"
	"net"
)

func main() {
	listenAddr := &net.UDPAddr{
		Port: 9000,
	}

	clientConn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		fmt.Println("error while connection via UDP:", err)
		return
	}
	defer clientConn.Close()

	fmt.Println("client up, waiting messages...")

	for {
		buffer := make([]byte, 1024)
		n, _, err := clientConn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("error while reading new message:", err)
			continue
		}

		fmt.Println("Local time:", string(buffer[:n]))
	}
}
