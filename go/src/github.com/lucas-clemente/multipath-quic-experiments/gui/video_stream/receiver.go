package video_stream

import (
    "log"
    "os/exec"
    "syscall"
    "bytes"
    "fmt"

    config "../config"
)

func ReceiverRoutine(quit chan bool, savePath string){
    log.Println("Save Path: ", savePath)
    c := exec.Command("go", "run", config.VID_RECEIVER_GO, savePath)
    c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
    
    var out bytes.Buffer
    var stderr bytes.Buffer
    c.Stdout = &out
    c.Stderr = &stderr

    err := c.Start()
    log.Println("stream_receiver.go started")

    for {
        select{
        case <- quit:
            pgid, err := syscall.Getpgid(c.Process.Pid)
            if err == nil {
                syscall.Kill(-pgid, 15)  // note the minus sign
            }

            err = c.Wait()
            if err != nil {
                log.Println(out.String())
                log.Println(err)
            }
            return
        default:
            if err != nil {
                log.Fatal("Error running stream_receiver.go")
                log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
            }
        }
    }
}