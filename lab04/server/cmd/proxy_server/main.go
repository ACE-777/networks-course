package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"server/internal"
)

const (
	port = 9000
)

func main() {
	internal.BuildBannedURL()

	proxy := &internal.ProxyServer{}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Printf("error while sw server starting: %v", err)

		return
	}

	defer listener.Close()

	log.Printf("Proxy server listening on port %v...", port)
	log.Fatal(http.Serve(listener, proxy))
}
