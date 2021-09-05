package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	utils "./utils"
	quic "github.com/lucas-clemente/quic-go"
)

// This script receives frames via TCP from the stream generator, and then transmits it to
// the peer using mpquic sockets.
// this script is acting as client for both tcp as well as mpquic server.

// config
const tcpServerAddr = "localhost:8002"
// const quicServerAddr = "192.168.43.148:4242"

func main() {

	quicServerAddr := os.Args[1] + ":4242"
	fmt.Println("quic server addr: ", quicServerAddr)
	// tcp connection
	tcpAddr, err := net.ResolveTCPAddr("tcp", tcpServerAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	// mpquic connection
	quicConfig := &quic.Config{
		CreatePaths: true,
	}

	sess, err := quic.DialAddr(quicServerAddr, &tls.Config{InsecureSkipVerify: true}, quicConfig)
	utils.HandleError(err)

	fmt.Println("session created: ", sess.RemoteAddr())

	stream, err := sess.OpenStream()
	utils.HandleError(err)

	defer stream.Close()
	defer stream.Close()

	fmt.Println("stream created...")

	// Infinite loop which takes frame from stream generator and then transmit it to the peer using
	// mpquic stream.
	for {
		// receive the frame size
		frame_size := make([]byte, 20)
		_, err = io.ReadFull(conn, frame_size)
		utils.HandleError(err)

		// fetch size value from the packet by removing trailing ':' symbols.
		size, _ := strconv.ParseInt(strings.Trim(string(frame_size), ":"), 10, 64)

		// terminate the process when an empty packet is received.
		if size == 0 {
			// send the reply size of zero to terminate the peer.
			stream.Write(frame_size)
			break
		}
		println("frame size: ", size)

		// receive the actual frame
		frame := make([]byte, size)
		_, err = io.ReadFull(conn, frame)
		utils.HandleError(err)

		// Send the frame size and frame using mpquic stream
		stream.Write(frame_size)
		stream.Write(frame)
	}
}
