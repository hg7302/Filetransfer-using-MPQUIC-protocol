package file_transfer

import (
    "log"
    "net"
    "os"

    "github.com/gotk3/gotk3/gtk"

    widgets "../widgets"
    utils "../utils"
)


func SetupSenderUI(win *gtk.Window) (*gtk.Grid){
    grid := widgets.GridNew(true, false, 5, 20)

    ipAddrEntry := widgets.EntryNew()

    pathLabel := widgets.LabelNew("<Path will appear hear>", true)
    fileChooserButton := widgets.ButtonNew("Click to Select File", func(){
        dialog := widgets.FileChooserDialogNew(win, gtk.FILE_CHOOSER_ACTION_OPEN)

        reply := dialog.Run()
        if reply == gtk.RESPONSE_OK {
            pathLabel.SetText(dialog.GetFilename())
        }
        dialog.Destroy()
    }, pathLabel)

    sendFileButton := widgets.ButtonNew("Send File", func (){
        button, err := grid.GetChildAt(0, 3)
        utils.HandleError(err)

        path, err := pathLabel.GetText()
        utils.HandleError(err)

        addr, err := ipAddrEntry.GetText()
        utils.HandleError(err)

        if err := net.ParseIP(addr); err != nil && utils.PathExists(path){
            // start client_multipath.go
            go SenderRoutine(button, path, addr)
            button.ToWidget().SetSensitive(false)
        } else {
            // set up a dialog
            log.Println("Wrong IP or path")
        }
    }, grid, pathLabel)


    grid.Attach(widgets.LabelNew("Server IP Address:", false), 0, 0, 2, 1)
    grid.Attach(ipAddrEntry, 2, 0, 2, 1)
    grid.Attach(widgets.LabelNew("Client IP Address:", false), 0, 1, 2, 1)
    grid.Attach(widgets.LabelNew(utils.GetOutboundIPAddr(), false), 2, 1, 2, 1)
    grid.Attach(widgets.LabelNew("Select file for transfer: ", false), 0, 2, 2, 1)
    grid.Attach(pathLabel, 2, 2, 2, 1)
    grid.Attach(sendFileButton, 0, 3, 2, 1)
    grid.Attach(fileChooserButton, 2, 3, 2, 1)

    return grid
}

func SetupReceiverUI(win *gtk.Window) (*gtk.Grid){
    grid := widgets.GridNew(true, false, 5, 20)

    receiverChannel := make(chan bool)
    currPath, err := os.Getwd()
    utils.HandleError(err)

    addrLabel := widgets.LabelNew(utils.GetOutboundIPAddr(), true)
    pathLabel := widgets.LabelNew(currPath, true)

    startButton := widgets.ButtonNew("Start", func (){
        button, err := grid.GetChildAt(2, 2)
        utils.HandleError(err)

        savePath, err := pathLabel.GetText()
        utils.HandleError(err)

        // start sender-multipath.go
        go ReceiverRoutine(receiverChannel, button, savePath)

        button.ToWidget().SetSensitive(false)
    }, grid, pathLabel)

    stopButton := widgets.ButtonNew("Stop", func (){
        button, err := grid.GetChildAt(2, 2)
        utils.HandleError(err)

        if !startButton.ToWidget().GetSensitive() {
            receiverChannel <- true
        }

        button.ToWidget().SetSensitive(true)
    }, grid)

    fileChooserButton := widgets.ButtonNew("Click to Select Folder", func(){
        dialog := widgets.FileChooserDialogNew(win, gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER)

        reply := dialog.Run()
        if reply == gtk.RESPONSE_OK {
            pathLabel.SetText(dialog.GetFilename())
        }
        dialog.Destroy()
    }, pathLabel)

    grid.Attach(widgets.LabelNew("Server IP Address:", false), 0, 0, 2, 1)
    grid.Attach(addrLabel, 2, 0, 2, 1)
    grid.Attach(widgets.LabelNew("Choose Path (for file-saving):", false), 0, 1, 2, 1)
    grid.Attach(pathLabel, 2, 1, 2, 1)
    grid.Attach(fileChooserButton, 0, 2, 2, 1)
    grid.Attach(startButton, 2, 2, 1, 1)
    grid.Attach(stopButton, 3, 2, 1, 1)
    return grid
}