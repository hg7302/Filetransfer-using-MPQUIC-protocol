package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	utils "./utils"
	quic "github.com/lucas-clemente/quic-go"
)

// This is a mpquic server which receives frames and saves them to jpeg files.

// quic server attaching to port 4242 on all the interfaces.
const quicServerAddr = "0.0.0.0:4242"

func main() {

	videoDir := os.Args[1]
	fmt.Println("Saving Video in: ", videoDir)

	quicConfig := &quic.Config{
		CreatePaths: true,
	}

	// initializing mpquic server
	fmt.Println("Attaching to: ", quicServerAddr)
	listener, err := quic.ListenAddr(quicServerAddr, utils.GenerateTLSConfig(), quicConfig)
	utils.HandleError(err)

	fmt.Println("Server started! Waiting for streams from client...")

	sess, err := listener.Accept()
	utils.HandleError(err)

	fmt.Println("session created: ", sess.RemoteAddr())

	stream, err := sess.AcceptStream()
	utils.HandleError(err)

	defer stream.Close()
	defer stream.Close()

	fmt.Println("stream created: ", stream.StreamID())

	frame_counter := 0

	// Infinite loop which receives frames from go sender peer and then saves them as jpeg files.
	for {
		frame_size := make([]byte, 20)
		_, err = io.ReadFull(stream, frame_size)

		size, _ := strconv.ParseInt(strings.Trim(string(frame_size), ":"), 10, 64)

		if size == 0 {
			break
		}

		fmt.Println("frame size: ", size)

		frame := make([]byte, size)

		_, err = io.ReadFull(stream, frame)

		jpeg_file, err := os.Create(videoDir + "/img" + strconv.Itoa(frame_counter) + ".jpg")
		utils.HandleError(err)
		frame_counter += 1

		jpeg_file.Write(frame)
		jpeg_file.Close()
	}
}
