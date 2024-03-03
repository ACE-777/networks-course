package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	portFromInput := flag.String("server_port", "8080", "configure port for ws server")
	flag.Parse()

	port := fmt.Sprintf(":%v", *portFromInput)
	listener, err := net.Listen("tcp", port)
	if err != nil {

		log.Printf("error while sw server starting: %v", err)

		return
	}

	defer listener.Close()

	log.Printf("Server is up on port%v", port)

	conn, err := listener.Accept()
	if err != nil {
		log.Printf("error while adding new connection: %v", err)

		return
	}

	defer conn.Close()

	log.Printf("connection estableshd")

	dataFromWsConnection, err := io.ReadAll(conn)
	if err != nil {
		log.Printf("error while reading data from user: %v", err)

		return
	}

	request := string(dataFromWsConnection)
	fmt.Println(request)

	lines := strings.Split(request, "\\")
	if len(lines) == 0 {
		log.Println("invalid request")

		return
	}

	parts := strings.Fields(lines[0])
	if len(parts) != 3 {
		log.Printf("invalid request")

		return
	}

	path := parts[1]

	path = strings.TrimLeft(path, "/")

	content, err := os.ReadFile(path)
	if err != nil {
		response := "HTTP/1.1 404 Not Found\r\n\r\nFile Not Found"
		conn.Write([]byte(response))

		return
	}

	response := "HTTP/1.1 200 OK\r\nContent-Length: " + fmt.Sprint(len(content)) + "\r\n\r\n" + string(content)
	_, err = conn.Write([]byte(response))
	if err != nil {
		log.Printf("%v", err)
	}
}
