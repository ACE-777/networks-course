package main

import (
	"bytes"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("can not up tcp server:%v", err)
	}
	log.Println("server up!")

	conn, err := listener.Accept()
	if err != nil {
		log.Fatalf("can not accept new connection:%v", err)
	}

	log.Println("server get new connection")

	var buf = make([]byte, 7)
	_, err = conn.Read(buf)
	if err != nil {
		log.Printf("can not get messages: %v", err)
	}

	log.Printf("server get new message with command:'%v'", string(buf))

	splitCommand := strings.Split(strings.TrimSpace(string(buf)), " ")
	cmd := exec.Command(splitCommand[0], splitCommand[1:]...)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("can not execute command:%v", err)
	}

	log.Printf("server successfully execute command")
	_, err = conn.Write(output.Bytes())
	if err != nil {
		log.Printf("can not send output of command: %v", err)
	}

	log.Println("successfully send result of executed command")
}
