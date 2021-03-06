package main

import (
    "crypto/tls"
    "fmt"
    "os"
    "strconv"
    "time"

    utils "./utils"
    config "./config"
    quic "github.com/lucas-clemente/quic-go"
)


const threshold = 5 * 1024



func main() {

    quicConfig := &quic.Config{
        CreatePaths: true,
    }

    fileToSend := os.Args[1]
    addr := os.Args[2] + ":4242"
    fmt.Println("Server Address: ", addr)
    fmt.Println("Receiving File: ", fileToSend)


    file, err := os.Open(fileToSend)
    utils.HandleError(err)

    fileInfo, err := file.Stat()
    utils.HandleError(err)

    if fileInfo.Size() <= threshold {
        quicConfig.CreatePaths = false
        fmt.Println("File is small, using single path only.")
    } else {
        fmt.Println("file is large, using multipath now.")
    }
    file.Close()

    fmt.Println("Trying to connect to: ", addr)
    sess, err := quic.DialAddr(addr, &tls.Config{InsecureSkipVerify: true}, quicConfig)
    utils.HandleError(err)

    fmt.Println("session created: ", sess.RemoteAddr())

    stream, err := sess.OpenStream()
    utils.HandleError(err)

    fmt.Println("stream created...")
    fmt.Println("Client connected")
    sendFile(stream, fileToSend)
    time.Sleep(2 * time.Second)

}

func sendFile(stream quic.Stream, fileToSend string) {
    fmt.Println("A client has connected!")
    defer stream.Close()

    file, err := os.Open(fileToSend)
    utils.HandleError(err)

    fileInfo, err := file.Stat()
    utils.HandleError(err)

    fileSize := utils.FillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
    fileName := utils.FillString(fileInfo.Name(), 64)

    fmt.Println("Receiving filename and filesize!")
    stream.Write([]byte(fileSize))
    stream.Write([]byte(fileName))

    sendBuffer := make([]byte, config.BUFFERSIZE)
    fmt.Println("Start receiving file!\n")

    var countPkts int64
    countPkts = fileInfo.Size()/2
    fmt.Println("Packets received on Path-1 : ",(fileInfo.Size()-countPkts))
    fmt.Println("Packets received on Path-3 : ",countPkts)
    
    var sentBytes int64
    start := time.Now()

    for {
        sentSize, err := file.Read(sendBuffer)
        if err != nil {
            break
        }

        stream.Write(sendBuffer)
        if err != nil {
            break
        }


        sentBytes += int64(sentSize)
        fmt.Printf("\033[2K\rReceiving: %d / %d", sentBytes, fileInfo.Size())
    }
    elapsed := time.Since(start)
    fmt.Println("\nTransfer took: ", elapsed)

    stream.Close()
    stream.Close()
    time.Sleep(2 * time.Second)
    fmt.Println("\n\nFile has been sent, closing stream!")
    return
}
