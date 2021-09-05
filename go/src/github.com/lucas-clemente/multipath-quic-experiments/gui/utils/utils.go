package utils

import (
    "log"
    "fmt"
    "os"
    "net"
)

func GetOutboundIPAddr() (string){
    var allIPAddr string

    ifaces, err := net.Interfaces()
    if err != nil {
        log.Fatal("Error retrieving IP Addrs:",err)
        return ""
    }
    for _, i := range ifaces {
        addrs, err := i.Addrs()
        if err != nil {
            log.Fatal("Error retrieving IP Addrs:",err)
            return ""
        }
        for _, addr := range addrs {
            ipnet, ok := addr.(*net.IPNet)

            if !ok{
                continue
            }
            ipv4 := ipnet.IP.To4()
            if ipv4 == nil || ipv4[0] == 127{
                continue
            }
            
            allIPAddr += ipv4.String() + " / "
        }
    }

    return allIPAddr
}

func HandleError(err error) {
    if err != nil {
        fmt.Println("Error: ", err)
        os.Exit(1)
    }
}

func PathExists(path string) (bool) {
    _, err := os.Stat(path)
    if err == nil { return true }
    if os.IsNotExist(err) { return false }
    return false
}
