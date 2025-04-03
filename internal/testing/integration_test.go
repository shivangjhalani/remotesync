package testing

import (
    "testing"
    "time"
    "github.com/shivangjhalani/remotesync/internal/network"
    "github.com/shivangjhalani/remotesync/internal/session"
    "github.com/shivangjhalani/remotesync/internal/security"
)

func TestFullWorkflow(t *testing.T) {
    // Start server
    server := network.NewServer()
    go server.Start(":8090")
    time.Sleep(time.Second) // Wait for server to start

    // Create two clients
    client1 := network.NewRemoteClient(network.ClientConfig{
        ServerHost: "localhost",
        ServerPort: "8090",
        Username:   "client1",
    })

    client2 := network.NewRemoteClient(network.ClientConfig{
        ServerHost: "localhost",
        ServerPort: "8090",
        Username:   "client2",
    })

    // Test connection
    if err := client1.Connect(); err != nil {
        t.Fatalf("Client 1 connection failed: %v", err)
    }
    if err := client2.Connect(); err != nil {
        t.Fatalf("Client 2 connection failed: %v", err)
    }

    // Test control request
    if err := client1.RequestControl(client2.GetID(), session.ModeSimultaneous); err != nil {
        t.Fatalf("Control request failed: %v", err)
    }

    // Test screen sharing
    if err := client2.StartScreenSharing(screen.CaptureConfig{
        FrameRate: 30,
        Quality:   75,
    }); err != nil {
        t.Fatalf("Screen sharing failed: %v", err)
    }
}