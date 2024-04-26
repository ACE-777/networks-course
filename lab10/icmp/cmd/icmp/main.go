package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	protocolICMP = 1
)

func main() {
	//if len(os.Args) != 2 {
	//	fmt.Println("Usage: ", os.Args[0], "hostname")
	//	os.Exit(1)
	//}

	//hostname := os.Args[1]
	hostname := "8.8.8.8"

	addr, err := net.ResolveIPAddr("ip4", hostname)
	if err != nil {
		fmt.Println("Error resolving address: ", err)
		os.Exit(1)
	}

	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Println("Error listening for ICMP packets: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Pinging", hostname, "...")
	for {
		sendTime := time.Now()
		seq := os.Getpid() & 0xffff
		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  seq,
				Data: []byte(""),
			},
		}
		msgBytes, err := msg.Marshal(nil)
		if err != nil {
			fmt.Println("Error marshaling ICMP message: ", err)
			continue
		}
		if _, err := conn.WriteTo(msgBytes, addr); err != nil {
			fmt.Println("Error sending ICMP message: ", err)
			continue
		}

		reply := make([]byte, 1500)
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		_, peer, err := conn.ReadFrom(reply)
		if err != nil {
			fmt.Println("Error reading ICMP reply: ", err)
			continue
		}
		recvTime := time.Now()

		recvMsg, err := icmp.ParseMessage(protocolICMP, reply)
		if err != nil {
			fmt.Println("Error parsing ICMP reply: ", err)
			continue
		}
		switch recvMsg.Type {
		case ipv4.ICMPTypeEchoReply:
			rtt := recvTime.Sub(sendTime)
			fmt.Printf("Reply from %s: seq=%d time=%v\n", peer, seq, rtt)
		default:
			fmt.Printf("Unexpected ICMP message: %+v\n", recvMsg)
		}

		select {
		case <-interrupt:
			fmt.Println("Ping interrupted.")
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
