package network_test

import (
    "testing"
    "remotesync/internal/network"
    "remotesync/internal/testutil"
)

func TestClientRegistration(t *testing.T) {
    server := network.NewServer()
    mockConn := &testutil.MockConn{}

    client, err := server.RegisterClient(mockConn)
    if err != nil {
        t.Errorf("Failed to register client: %v", err)
    }

    if client.ID == "" {
        t.Error("Client ID not generated")
    }

    if len(server.GetClientList()) != 1 {
        t.Error("Client not added to list")
    }
}

func TestMaxClients(t *testing.T) {
    server := network.NewServer()
    
    // Register maximum allowed clients
    for i := 0; i < 2; i++ {
        _, err := server.RegisterClient(&testutil.MockConn{})
        if err != nil {
            t.Errorf("Failed to register client %d: %v", i, err)
        }
    }

    // Try registering one more client
    _, err := server.RegisterClient(&testutil.MockConn{})
    if err == nil {
        t.Error("Server allowed more than maximum clients")
    }
}