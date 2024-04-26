package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	protocolICMP = 1
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "hostname")
		os.Exit(1)
	}

	hostname := os.Args[1]

	addr, err := net.ResolveIPAddr("ip4", hostname)
	if err != nil {
		log.Printf("Error resolving address: %v", err)
	}

	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Printf("Error listening for ICMP packets: %v", err)
	}

	defer conn.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var (
		rttMin, rttMax, rttTotal     time.Duration
		packetsSent, packetsReceived int
	)

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

		packetsSent++

		reply := make([]byte, 1500)
		if err = conn.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
			fmt.Printf("Can not set read deadline timeout: %v", err)
		}

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
			packetsReceived++

			if packetsReceived == 1 || rtt < rttMin {
				rttMin = rtt
			}

			if rtt > rttMax {
				rttMax = rtt
			}

			rttTotal += rtt
		default:
			fmt.Printf("Unexpected ICMP message: %+v\n", recvMsg)
		}

		printStatistics(packetsSent, packetsReceived, rttMin, rttMax, rttTotal)

		select {
		case <-interrupt:
			fmt.Println("Ping interrupted.")

			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func printStatistics(packetsSent, packetsReceived int, rttMin, rttMax, rttTotal time.Duration) {
	packetLoss := float64(packetsSent-packetsReceived) / float64(packetsSent) * 100
	if packetsReceived > 0 {
		rttAvg := rttTotal / time.Duration(packetsReceived)
		fmt.Printf("Ping statistics: Sent=%d, Received=%d, Loss=%.2f%%, MinRTT=%v, MaxRTT=%v, AvgRTT=%v\n",
			packetsSent, packetsReceived, packetLoss, rttMin, rttMax, rttAvg)
	} else {
		fmt.Printf("Ping statistics: Sent=%d, Received=%d, Loss=%.2f%%\n",
			packetsSent, packetsReceived, packetLoss)
	}
}
