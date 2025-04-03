package protocol_test

import (
    "testing"
    "github.com/shivangjhalani/remotesync/internal/protocol"
)

func TestMessageSerialization(t *testing.T) {
    original := protocol.NewMessage(protocol.TypePing, map[string]any{
        "test": "data",
    })

    data, err := original.Encode()
    if err != nil {
        t.Errorf("Failed to encode message: %v", err)
    }

    decoded, err := protocol.DecodeMessage(data)
    if err != nil {
        t.Errorf("Failed to decode message: %v", err)
    }

    if decoded.Type != original.Type {
        t.Errorf("Message type mismatch: got %v, want %v", 
                 decoded.Type, original.Type)
    }

    if decoded.Payload["test"] != "data" {
        t.Error("Payload data mismatch")
    }
}