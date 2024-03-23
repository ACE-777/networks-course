package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	smtpServer := "127.0.0.1"
	smtpPort := "2500"

	from := "st072182@student.spbu.ru"
	to := "st072182@student.spbu.ru"
	body := "Hi MKN!"

	conn, err := net.Dial("tcp", smtpServer+":"+smtpPort)
	if err != nil {
		fmt.Println("Error connecting to SMTP server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	greeting, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading greeting from server:", err)
		return
	}
	fmt.Println("Server greeting:", greeting)

	fmt.Fprintf(conn, "EHLO example.com\r\n")
	respEHLO, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error sending EHLO command:", err)
		return
	}
	fmt.Println("EHLO response:", respEHLO)

	fmt.Fprintf(conn, "MAIL FROM:<%s>\r\n", from)
	respMAILFROM, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error sending MAIL FROM command:", err)
		return
	}
	fmt.Println("MAIL FROM response:", respMAILFROM)

	fmt.Fprintf(conn, "RCPT TO:<%s>\r\n", to)
	respRCPTTO, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error sending RCPT TO command:", err)
		return
	}
	fmt.Println("RCPT TO response:", respRCPTTO)

	fmt.Fprintf(conn, "DATA\r\n")

	respDATA, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error sending DATA command:", err)
		return
	}
	fmt.Println("DATA response:", respDATA)

	fmt.Fprintf(conn, "Subject: Test\r\n")
	fmt.Fprintf(conn, "From: <%s>\r\n", from)
	fmt.Fprintf(conn, "To: <%s>\r\n", to)
	fmt.Fprintf(conn, "\r\n\r\n")

	fmt.Fprintf(conn, "%s.\r\n", body)

	fmt.Fprintf(conn, "\r\n.\r\n")
	respEnd, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error sending end of message:", err)
		return
	}
	fmt.Println("End of message response:", respEnd)

	fmt.Fprintf(conn, "QUIT\r\n")
	respQUIT, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error sending QUIT command:", err)
		return
	}
	fmt.Println("QUIT response:", respQUIT)

	fmt.Println("Email sent successfully!")
}
