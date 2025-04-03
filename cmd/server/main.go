package main

import (
    "fmt"
    "net"
    "github.com/shivangjhalani/remotesync/internal/network"
    "github.com/shivangjhalani/remotesync/internal/logger"
)

const (
    HOST = "localhost"
    PORT = "8080"
)

var server *network.Server

func main() {
    server = network.NewServer()
    
    listener, err := net.Listen("tcp", HOST+":"+PORT)
    if err != nil {
        logger.ErrorLogger.Fatalf("Error starting server: %v", err)
    }
    defer listener.Close()

    logger.InfoLogger.Printf("Server listening on %s:%s", HOST, PORT)

    for {
        conn, err := listener.Accept()
        if (err != nil) {
            logger.ErrorLogger.Printf("Error accepting connection: %v", err)
            continue
        }

        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    client, err := server.RegisterClient(conn)
    if err != nil {
        logger.ErrorLogger.Printf("Failed to register client: %v", err)
        conn.Close()
        return
    }

    logger.InfoLogger.Printf("New client registered: %s", client.ID)
    defer server.RemoveClient(client.ID)

    // TODO: Handle client communication
}