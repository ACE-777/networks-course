package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	handler_request "lab_3/internal"
)

func main() {
	port := flag.String("server_port", "9000", "ws server port")
	concurrencyLevel := flag.Int("concurrency_level", 5, "max active connection")

	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", *port))
	if err != nil {
		log.Printf("error while sw server starting: %v", err)

		return
	}

	defer listener.Close()

	log.Printf("Server is up on port%v", *port)

	handler_request.Sem = make(chan struct{}, *concurrencyLevel)

	for {
		conn, err := listener.Accept()
		log.Printf("accept connection")
		if err != nil {
			log.Printf("error while adding new connection: %v", err)

			continue
		}

		handler_request.Sem <- struct{}{}

		go handler_request.Handler(conn, handler_request.Sem)
	}
}
