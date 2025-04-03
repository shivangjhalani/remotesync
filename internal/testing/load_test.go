package testing

import (
    "testing"
    "sync"
    "remotesync/internal/network"
)

func TestMultipleConnections(t *testing.T) {
    server := network.NewServer()
    go server.Start(":8091")

    var wg sync.WaitGroup
    errorChan := make(chan error, 10)

    // Try connecting multiple clients simultaneously
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            client := network.NewRemoteClient(network.ClientConfig{
                ServerHost: "localhost",
                ServerPort: "8091",
                Username:  fmt.Sprintf("test-client-%d", id),
            })
            if err := client.Connect(); err != nil {
                errorChan <- err
            }
        }(i)
    }

    wg.Wait()
    close(errorChan)

    for err := range errorChan {
        t.Errorf("Connection error: %v", err)
    }
}