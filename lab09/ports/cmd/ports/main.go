package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Input: ./main.go <IP> <start port> <end port>")
		return
	}

	ip := os.Args[1]
	startPort, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Can't take startport:%v", err)
	}

	endPort, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatalf("Can't take endport:%v", err)
	}

	fmt.Printf("Scan ports on %s from %d to %d\n", ip, startPort, endPort)
	for port := startPort; port <= endPort; port++ {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
		if err != nil {
			fmt.Printf("Port %d closed\n", port)
		} else {
			defer conn.Close()
			fmt.Printf("Port %d open\n", port)
		}
	}
}
