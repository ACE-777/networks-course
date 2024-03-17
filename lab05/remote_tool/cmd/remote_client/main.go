package main

import (
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		log.Fatalf("can not connect to tcp server:%v", err)
	}

	log.Println("connection established")

	_, err = conn.Write([]byte("echo hi"))
	if err != nil {
		log.Printf("can not send messages: %v", err)
	}

	log.Println("successfully send message")

	{
		var buf2 = make([]byte, 1024)
		_, err = conn.Read(buf2)
		if err != nil {
			log.Printf("can not get messages: %v", err)
		}

		log.Printf("result of command: %v", string(buf2))
	}
}
