package performance

import (
    "testing"
    "time"
    "github.com/shivangjhalani/remotesync/internal/screen"
    "github.com/shivangjhalani/remotesync/internal/network"
    "github.com/shivangjhalani/remotesync/internal/protocol"
    "github.com/shivangjhalani/remotesync/internal/testutil"
)

func BenchmarkScreenCapture(b *testing.B) {
    capture := screen.NewScreenCapture(30, 75)
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        frame := <-capture.frameChan
        if frame == nil {
            b.Fatal("Failed to capture frame")
        }
    }
}

func BenchmarkFrameCompression(b *testing.B) {
    capture := screen.NewScreenCapture(30, 75)
    frame := <-capture.frameChan
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := screen.CompressFrame(frame.Image, 75)
        if err != nil {
            b.Fatal("Frame compression failed:", err)
        }
    }
}

func BenchmarkMessageTransmission(b *testing.B) {
    mockConn := &testutil.MockConn{}
    client := network.NewRemoteClient(network.ClientConfig{})
    client.SetConnection(mockConn)
    
    msg := protocol.NewMessage(protocol.TypePing, map[string]any{
        "timestamp": time.Now(),
    })
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        err := client.SendMessage(msg)
        if err != nil {
            b.Fatal("Message transmission failed:", err)
        }
    }
}

func BenchmarkConcurrentInputHandling(b *testing.B) {
    server := network.NewServer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            mockConn := &testutil.MockConn{}
            _, err := server.RegisterClient(mockConn)
            if err != nil {
                b.Fatal("Client registration failed:", err)
            }
        }
    })
}