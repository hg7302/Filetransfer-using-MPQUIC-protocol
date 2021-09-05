package video_stream

import (
    "fmt"
    "time"
    "os"
    "log"
    "net"

    "github.com/gotk3/gotk3/gtk"

    // config "../config"
    utils "../utils"
    widgets "../widgets"
)


func SetupVideo(quit chan bool, img *gtk.Image, dir string){
    time.Sleep(3 * time.Second)

    counter := 0
    path := dir + "/img%d.jpg"
    
    for {
        select{
        case <- quit:
            img.Clear()
            return
        default:
            currentPath := fmt.Sprintf(path, counter)

            if utils.PathExists(currentPath) {
                img.SetFromFile(currentPath)
                // time.Sleep(200 * time.Millisecond)
                counter += 1
                img.Show()
            }                
        }
    }
}


func SetupReceiverUI(win *gtk.Window) (*gtk.Grid) {
    grid := widgets.GridNew(true, false, 5, 20)
    img := widgets.ImageNew()
    
    streamReceiverChannel := make(chan bool)
    videoGUIChannel := make(chan bool)
    

    path, err := os.Getwd()
    utils.HandleError(err)
    pathLabel := widgets.LabelNew(path, true)
    
    fileChooserButton := widgets.ButtonNew("Click to Select Folder", func(){
        dialog := widgets.FileChooserDialogNew(win, gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER)

        reply := dialog.Run()
        if reply == gtk.RESPONSE_OK {
            pathLabel.SetText(dialog.GetFilename())
        }
        dialog.Destroy()
    }, pathLabel)

    startButton := widgets.ButtonNew("Start", func (){
        button, err := grid.GetChildAt(2, 2)
        utils.HandleError(err)

        path, err := pathLabel.GetText()
        utils.HandleError(err)
        
        // start stream_receiver.go
        go ReceiverRoutine(streamReceiverChannel, path)
        // start videoGUIChannel
        go SetupVideo(videoGUIChannel, img, path)

        button.ToWidget().SetSensitive(false)
    }, grid, pathLabel, img)


    stopButton := widgets.ButtonNew("Stop", func (){
        button, err := grid.GetChildAt(2, 2)
        utils.HandleError(err)

        if sensitive := button.ToWidget().GetSensitive(); !sensitive{
            // stop stream_receiver.go
            streamReceiverChannel <- true
            // stop video stream on GUI
            videoGUIChannel <- true
        }

        button.ToWidget().SetSensitive(true)
    }, grid)

    grid.Attach(widgets.LabelNew("Server IP Address: ", false), 0, 0, 2, 1)
    grid.Attach(widgets.LabelNew(utils.GetOutboundIPAddr(), true), 2, 0, 2, 1)
    grid.Attach(widgets.LabelNew("Folder for images:", true), 0, 1, 2, 1)
    grid.Attach(pathLabel, 2, 1, 2, 1) 
    grid.Attach(fileChooserButton, 0, 2, 2, 1)
    grid.Attach(startButton, 2, 2, 1, 1)
    grid.Attach(stopButton, 3, 2, 1, 1)
    grid.Attach(img, 0, 5, 4, 4)
    return grid

}

func SetupSenderUI(win *gtk.Window) (*gtk.Grid) {
    grid := widgets.GridNew(true, false, 5, 20)
    img := widgets.ImageNew()

    ipAddrEntry := widgets.EntryNew()

    streamGeneratorChannel := make(chan bool)
    streamSenderChannel := make(chan bool)
    // videoGUIChannel := make(chan bool)

    startButton := widgets.ButtonNew("Start", func(){
        button, err := grid.GetChildAt(1, 1)
        utils.HandleError(err)

        addr, err := ipAddrEntry.GetText()
        utils.HandleError(err)

        if err := net.ParseIP(addr); err != nil{
            // start stream_generator.py
            go GeneratorRoutine(streamGeneratorChannel)
            // start stream_sender.go
            go SenderRoutine(streamSenderChannel, addr)
            // start video stream on GUI
            // go streamVideoOnGUI(videoGUIChannel, img)
            button.ToWidget().SetSensitive(false)
        } else {
            log.Println("Wrong IP")
        }
    }, streamGeneratorChannel, streamSenderChannel, ipAddrEntry)

    stopButton := widgets.ButtonNew("Stop", func (){
        button, err := grid.GetChildAt(1, 1)
        utils.HandleError(err)
        
        if sensitive := button.ToWidget().GetSensitive(); !sensitive{
            // stop stream_generator.py
            streamSenderChannel <- true
            // stop stream_sender.go
            streamGeneratorChannel <- true
            // stop video stream on GUI
            // videoGUIChannel <- true
        }
        button.ToWidget().SetSensitive(true)
    }, streamGeneratorChannel, streamSenderChannel, grid)


    grid.Attach(widgets.LabelNew("Server IP Address:", true), 0, 0, 2, 1)
    grid.Attach(ipAddrEntry, 2, 0, 2, 1)
    grid.Attach(startButton, 1, 1, 1, 1)
    grid.Attach(stopButton, 2, 1, 1, 1)
    grid.Attach(img, 0, 3, 4, 4)
    return grid
}
