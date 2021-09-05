package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	utils "./utils"
	config "./config"
	quic "github.com/lucas-clemente/quic-go"
)

const addr = "0.0.0.0:" + config.PORT

func main() {

	savePath := os.Args[1]
	fmt.Println("Saving file to: ", savePath)

	quicConfig := &quic.Config{
		CreatePaths: true,
	}

	fmt.Println("Attaching to: ", addr)
	listener, err := quic.ListenAddr(addr, utils.GenerateTLSConfig(), quicConfig)
	utils.HandleError(err)

	fmt.Println("Server started! Waiting for streams from client...")

	sess, err := listener.Accept()
	utils.HandleError(err)

	fmt.Println("session created: ", sess.RemoteAddr())

	stream, err := sess.AcceptStream()
	utils.HandleError(err)

	fmt.Println("stream created: ", stream.StreamID())
	fmt.Println("Priority of stream: %+v", stream.Priority())

	defer stream.Close()
	fmt.Println("Connected to server, start sending the file name and file size")
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	stream.Read(bufferFileSize)
	fmt.Println("stream created: ", stream.StreamID())
	fmt.Println("Priority of stream: %+v", stream.Priority())
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	fmt.Println("file size sent: ", fileSize)

	stream.Read(bufferFileName)
	fmt.Println("stream created: ", stream.StreamID())
	fmt.Println("Priority of stream: %+v", stream.Priority())
	fileName := strings.Trim(string(bufferFileName), ":")

	fmt.Println("file name sent: ", fileName)

	newFile, err := os.Create(savePath + "/" + fileName)
	utils.HandleError(err)

	defer newFile.Close()
	var receivedBytes int64
	start := time.Now()

	for {
		if (fileSize - receivedBytes) < config.BUFFERSIZE {
			recv, err := io.CopyN(newFile, stream, (fileSize - receivedBytes))
			utils.HandleError(err)

			stream.Read(make([]byte, (receivedBytes + config.BUFFERSIZE) - fileSize))
			receivedBytes += recv
			fmt.Printf("\033[2K\rSent: %d / %d", receivedBytes, fileSize)

			break
		}
		_, err := io.CopyN(newFile, stream, config.BUFFERSIZE)
		utils.HandleError(err)

		receivedBytes += config.BUFFERSIZE

		fmt.Printf("\033[2K\rSent: %d / %d", receivedBytes, fileSize)
	}
	elapsed := time.Since(start)
	fmt.Println("\nTransfer took: ", elapsed)

	time.Sleep(2 * time.Second)
	stream.Close()
	stream.Close()
	fmt.Println("\n\nReceived file completely!")
}
