package main

import (
    "log"
    "fmt"

    "github.com/gotk3/gotk3/gtk"
    
    // config "./config"
    file_transfer "./file_transfer"
    // utils "./utils"
    video_stream "./video_stream"
    widgets "./widgets"
)

var (
    role = "server"
    close = false
    win *gtk.Window
)

func addClientSide(win *gtk.Window) {
    stackSwitcher := widgets.StackSwitcherNew()  
    stack := widgets.StackNew()

    gridFileTransfer := file_transfer.SetupSenderUI(win)
    gridVideoStream := video_stream.SetupSenderUI(win)
    stack.AddTitled(gridFileTransfer, "Page1", "File Transfer")
    stack.AddTitled(gridVideoStream, "Page2", "Video Stream")
    stackSwitcher.SetStack(stack)

    box := widgets.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    box.PackStart(stackSwitcher, false, false, 0)
    box.PackStart(stack, true, true, 0)

    win.Add(box)
}


func addServerSide(win *gtk.Window){
    stackSwitcher := widgets.StackSwitcherNew()
    stack := widgets.StackNew()

    gridFileTransfer := file_transfer.SetupReceiverUI(win)
    gridVideoStream := video_stream.SetupReceiverUI(win)
    stack.AddTitled(gridFileTransfer, "Page1", "File Transfer")
    stack.AddTitled(gridVideoStream, "Page2", "Video Stream")
    stackSwitcher.SetStack(stack)

    box := widgets.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    box.PackStart(stackSwitcher, false, false, 0)
    box.PackStart(stack, true, true, 0)

    win.Add(box)
}


func setupDialog() (){
    dialog := widgets.DialogNew("MPQUIC Experiment", 300, 150)

    dialog.AddButton("OK", gtk.RESPONSE_OK)
    dialog.AddButton("Cancel", gtk.RESPONSE_CLOSE)

    contentArea, err := dialog.GetContentArea()
    if err != nil {
        log.Fatal("Unable to fetch contentArea: ", err)
    }

    box := widgets.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    clientButton := widgets.RadioButtonNew(nil, "Server", func(){
        role = "server"
    })
    serverButton := widgets.RadioButtonNew(clientButton, "Client", func(){
        role = "client"
    })
    box.PackStart(clientButton, false, false, 0)
    box.PackStart(serverButton, false, false, 0)


    contentArea.PackStart(widgets.LabelNew("Which role do you want to start?", false), false, false, 0)
    contentArea.PackStart(box, false, false, 0)
    dialog.ShowAll()

    reply := dialog.Run()
    if reply == gtk.RESPONSE_OK {
        fmt.Println("OK")
    } else {
        close = true
    }
    dialog.Destroy()
}

func main(){
    gtk.Init(nil)

    setupDialog()

    log.Printf("Selected Role: %s", role)

    if role == "client"{
        win = widgets.WindowNew("Client", 800, 200)
        addClientSide(win)
    } else {
        win = widgets.WindowNew("Server", 800, 600)
        addServerSide(win)
    }
    
    if !close {
        win.ShowAll()
        gtk.Main()
    }

}