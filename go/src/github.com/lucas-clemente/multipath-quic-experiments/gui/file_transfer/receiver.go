package file_transfer

import (
    "log"
    "bytes"
    "os/exec"
    "syscall"
    "fmt"

    "github.com/gotk3/gotk3/gtk"

    config "../config"
)

func ReceiverRoutine (quit chan bool, button gtk.IWidget, path string){
    c := exec.Command("go", "run", config.FILE_SERVER_GO, path)
    c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

    var out bytes.Buffer
    var stderr bytes.Buffer

    c.Stdout = &out
    c.Stderr = &stderr

    err := c.Start()
    if err != nil{
        log.Println("server-multipath.go did not start correctly")
        log.Println("Error code :", fmt.Sprint(err), ": ", stderr.String())
        return
    } else {
        log.Println("server-multipath.go started")
    }


    done := make(chan bool)
    go (func() {
        c.Wait()
        close(done)
    })()

    select{
    case <- done:
        log.Println("server-multipath.go finished normally")
        button.ToWidget().SetSensitive(true)
        return
    case <- quit:
        log.Println("Killing server-multipath.go")
        pgid, err := syscall.Getpgid(c.Process.Pid)
        if err == nil {
            syscall.Kill(-pgid, 15)  // note the minus sign
        }

        err = c.Wait()
        if err != nil {
            log.Println(err)
        }
        button.ToWidget().SetSensitive(true)
        return
    }
}
