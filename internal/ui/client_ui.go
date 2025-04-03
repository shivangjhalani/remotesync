package ui

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "remotesync/internal/network"
    "remotesync/internal/logger"
)

type ClientUI struct {
    window     fyne.Window
    client     *network.RemoteClient
    clientList *widget.List
    clients    []string
}

func NewClientUI(client *network.RemoteClient) *ClientUI {
    myApp := app.New()
    window := myApp.NewWindow("RemoteSync Client")
    
    return &ClientUI{
        window: window,
        client: client,
    }
}

func (ui *ClientUI) Setup() {
    // Connection settings
    hostEntry := widget.NewEntry()
    hostEntry.SetPlaceHolder("Server Host")
    portEntry := widget.NewEntry()
    portEntry.SetPlaceHolder("Server Port")

    connectBtn := widget.NewButton("Connect", func() {
        config := network.ClientConfig{
            ServerHost: hostEntry.Text,
            ServerPort: portEntry.Text,
        }
        if err := ui.client.Connect(); err != nil {
            logger.ErrorLogger.Printf("Connection failed: %v", err)
        }
    })

    // Available clients list
    ui.clientList = widget.NewList(
        func() int { return len(ui.clients) },
        func() fyne.CanvasObject { return widget.NewLabel("Client") },
        func(i widget.ListItemID, o fyne.CanvasObject) {
            o.(*widget.Label).SetText(ui.clients[i])
        },
    )

    // Control mode selection
    modeSelect := widget.NewSelect([]string{"Exclusive", "Simultaneous"}, func(mode string) {
        // Handle mode selection
    })

    // Settings panel
    qualitySlider := widget.NewSlider(0, 100)
    fpsSlider := widget.NewSlider(1, 60)

    settingsContainer := container.NewVBox(
        widget.NewLabel("Quality"),
        qualitySlider,
        widget.NewLabel("FPS"),
        fpsSlider,
    )

    // Main layout
    content := container.NewHSplit(
        container.NewVBox(
            hostEntry,
            portEntry,
            connectBtn,
            ui.clientList,
        ),
        container.NewVBox(
            modeSelect,
            settingsContainer,
        ),
    )

    ui.window.SetContent(content)
}

func (ui *ClientUI) Show() {
    ui.window.ShowAndRun()
}