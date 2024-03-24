package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	internal "lab05/internal"

	"github.com/jlaffaye/ftp"
)

func main() {
	c, err := ftp.Dial("127.0.0.1:21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully dialed FTP server")

	defer func() {
		if err := c.Quit(); err != nil {
			log.Fatal(err)
		}
	}()

	err = c.Login("TestUser", "")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to FTP server")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("ftp> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		input = strings.TrimSpace(input)

		if input == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		fields := strings.Fields(input)
		command := fields[0]

		switch command {
		case "ls":
			dirList, err := c.List(".")
			if err != nil {
				log.Println("Error listing directory:", err)
			}

			fmt.Println("Files and directories on FTP server:")
			for _, item := range dirList {
				fmt.Println(item.Name)
			}

		case "get":
			if len(fields) != 3 {
				fmt.Println("Usage: get <remote_file> <local_file>")
				continue
			}

			remotePath := fields[1]
			localPath := fields[2]
			internal.DownloadFile(c, remotePath, localPath)

		case "put":
			if len(fields) != 3 {
				fmt.Println("Usage: put <local_file> <remote_file>")
				continue
			}

			localPath := fields[1]
			remotePath := fields[2]
			internal.UploadFile(c, localPath, remotePath)

		default:
			fmt.Println("Unknown command:", command)
		}
	}
}
