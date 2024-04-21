package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

func main() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("Can get ip adress")
	}

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			fmt.Println("IP:", ipnet.IP)
			mask := ipnet.IP.DefaultMask()
			fmt.Println("Mask:")
			for _, b := range mask {
				fmt.Print(".", strconv.Itoa(int(b)))
			}

			break
		}
	}
}
