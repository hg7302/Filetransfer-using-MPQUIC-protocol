package video_stream

import (
    "os/exec"
    "log"
    "bytes"
    "fmt"

    config "../config"
)


func GeneratorRoutine(quit chan bool){
    c := exec.Command(config.PYTHON, config.VID_GENERATOR_PY)
    var out bytes.Buffer
    var stderr bytes.Buffer
    
    c.Stdout = &out
    c.Stderr = &stderr
    
    err := c.Start()
    log.Println("stream_generator.py started")

    for {
        select{
        case <- quit:
            if err := c.Process.Kill(); err != nil{
                log.Println("Error occurred while killing stream_generator.py")
            } else {
                log.Println(out.String())
                log.Println("Killed stream_generator.py process")
            }
            return
        default:
            if err != nil {
                log.Fatal("Error running stream_generator.py")
                log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
            } 
        }
    }
}

func SenderRoutine(quit chan bool, addr string){
    c := exec.Command("go","run", config.VID_SENDER_GO, addr)
    var out bytes.Buffer
    var stderr bytes.Buffer
    
    c.Stdout = &out
    c.Stderr = &stderr
    
    err := c.Start()
    log.Println("stream_sender.go started")
    
    for {
        select{
        case <- quit:
            if err := c.Process.Kill(); err != nil{
                log.Println("Error occurred while killing stream_sender.go")
            } else {
                log.Println(out.String())
                log.Println("Killed stream_sender.go process")
            }
            return
        default:
            if err != nil {
                log.Fatal("Error running stream_sender.go")
                log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
            }
        }
    }
}
