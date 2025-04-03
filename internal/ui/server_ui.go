package ui

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "remotesync/internal/network"
)

type ServerUI struct {
    window   fyne.Window
    server   *network.Server
    clients  *widget.Table
    sessions *widget.List
}

func NewServerUI(server *network.Server) *ServerUI {
    myApp := app.New()
    window := myApp.NewWindow("RemoteSync Server")
    
    return &ServerUI{
        window: window,
        server: server,
    }
}

func (ui *ServerUI) Setup() {
    // Server status
    statusLabel := widget.NewLabel("Server Status: Running")
    
    // Connected clients table
    ui.clients = widget.NewTable(
        func() (int, int) { return len(ui.server.GetClientList()), 3 },
        func() fyne.CanvasObject { return widget.NewLabel("") },
        func(i widget.TableCellID, o fyne.CanvasObject) {
            client := ui.server.GetClientList()[i.Row]
            label := o.(*widget.Label)
            switch i.Col {
            case 0:
                label.SetText(client.ID)
            case 1:
                label.SetText(client.Name)
            case 2:
                label.SetText(client.Conn.RemoteAddr().String())
            }
        },
    )

    // Configuration
    portEntry := widget.NewEntry()
    portEntry.SetPlaceHolder("Server Port")
    maxClientsEntry := widget.NewEntry()
    maxClientsEntry.SetText("2")

    configContainer := container.NewVBox(
        widget.NewLabel("Configuration"),
        portEntry,
        widget.NewLabel("Max Clients"),
        maxClientsEntry,
    )

    // Main layout
    content := container.NewHSplit(
        container.NewVBox(
            statusLabel,
            ui.clients,
        ),
        container.NewVBox(
            configContainer,
            ui.sessions,
        ),
    )

    ui.window.SetContent(content)
}

func (ui *ServerUI) Show() {
    ui.window.ShowAndRun()
}