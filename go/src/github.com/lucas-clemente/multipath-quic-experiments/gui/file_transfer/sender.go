package file_transfer

import (
    "log"
    "os/exec"
    "bytes"
    "fmt"

    "github.com/gotk3/gotk3/gtk"

    config "../config"
)

func SenderRoutine(button gtk.IWidget, path string, addr string) {
	log.Println("File to be sent : ", path)
    c := exec.Command("go", "run", config.FILE_CLIENT_GO, path, addr)
    var out bytes.Buffer
    var stderr bytes.Buffer

    c.Stdout = &out
    c.Stderr = &stderr

    err := c.Start()
    if err != nil {
        log.Fatal("Error while running client-multipath.go")
    } else {
	    log.Println("client-multipath.go started")
	}

    if err := c.Wait(); err != nil {
        log.Println("client-multipath.go did not finished correctly")
        log.Println("Error code :", fmt.Sprint(err), ": ", stderr.String())
    } else {
        log.Println("client-multipath.go finished")
        log.Println(out.String())
    }
    button.ToWidget().SetSensitive(true)

}