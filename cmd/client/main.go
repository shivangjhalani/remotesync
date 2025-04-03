package main

import (
    "flag"
    "time"
    "github.com/shivangjhalani/remotesync/internal/network"
    "github.com/shivangjhalani/remotesync/internal/logger"
)

var (
    serverHost = flag.String("host", "localhost", "Server host address")
    serverPort = flag.String("port", "8080", "Server port")
)

func main() {
    flag.Parse()

    config := network.ClientConfig{
        ServerHost: *serverHost,
        ServerPort: *serverPort,
        RetryDelay: 5 * time.Second,
        MaxRetries: 5,
    }

    client := network.NewRemoteClient(config)
    
    logger.InfoLogger.Printf("Connecting to server at %s:%s", *serverHost, *serverPort)
    
    if err := client.Connect(); err != nil {
        logger.ErrorLogger.Fatalf("Failed to connect to server: %v", err)
    }

    // Keep the client running
    select {}
}