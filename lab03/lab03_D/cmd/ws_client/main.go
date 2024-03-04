package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	internal "lab_3/internal"
)

func main() {
	host := flag.String("server_host", "127.0.0.1", "address of host server")
	port := flag.String("server_port", "9000", "ws server port")
	file := flag.String("file", "files/example_file.txt", "path to request file")

	flag.Parse()

	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", *host, *port))
	if err != nil {
		log.Printf("can not cnnect to host %v with port %v", *host, *port)
	}

	defer conn.Close()
	log.Printf("succesfully connected to server on host: %v with port: %v via tcp", *host, *port)
	request := fmt.Sprintf("GET /%s HTTP/1.1\r\nHost: %s\r\n\r\n", *file, fmt.Sprintf("%v:%v", *host, *port))
	_, err = conn.Write([]byte(request))
	if err != nil {
		log.Printf("%v", err)
	}

	log.Printf("succesfully send new message")

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		log.Printf("error while reading data from user: %v", err)

		return
	}

	response := string(buf)
	fmt.Println(response)

	<-internal.Sem
}
