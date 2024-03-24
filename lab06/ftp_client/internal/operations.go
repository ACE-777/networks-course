package internal

import (
	"github.com/jlaffaye/ftp"
	"io"
	"log"
	"os"
)

func DownloadFile(c *ftp.ServerConn, remotePath, localPath string) {
	downloadedFile, err := os.Create(localPath)
	if err != nil {
		log.Println("Error creating local file:", err)
		return
	}

	defer downloadedFile.Close()
	retr, err := c.Retr(remotePath)
	if err != nil {
		return
	}

	_, err = io.Copy(downloadedFile, retr)
	if err != nil {
		log.Println("Error copy file to local via downloading:", err)
	}

	if err != nil {
		log.Println("Error downloading file:", err)
		return
	}

	log.Println("File", remotePath, "downloaded and saved to", localPath)
}

func UploadFile(c *ftp.ServerConn, localPath, remotePath string) {
	file, err := os.Open(localPath)
	if err != nil {
		log.Println("Error opening local file:", err)
		return
	}

	defer file.Close()

	err = c.Stor(remotePath, file)
	if err != nil {
		log.Println("Error uploading file:", err)
		return
	}

	log.Println("File", localPath, "uploaded to", remotePath)
}
