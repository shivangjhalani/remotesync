package protocol

import (
    "encoding/json"
    "time"
)

type MessageType string

const (
    TypeIdentify      MessageType = "identify"
    TypePing         MessageType = "ping"
    TypePong         MessageType = "pong"
    TypeClientList   MessageType = "client_list"
    TypeControlReq   MessageType = "control_request"
    TypeControlRes   MessageType = "control_response"
    TypeMouseEvent   MessageType = "mouse_event"
    TypeKeyEvent     MessageType = "key_event"
    TypeScreenData   MessageType = "screen_data"
    TypeError        MessageType = "error"
)

// Add documentation comments
type Message struct {
    Type      MessageType     `json:"type"`      // Type of the message
    Timestamp time.Time       `json:"timestamp"` // When the message was created
    Payload   map[string]any  `json:"payload"`   // Message data
    SenderID  string         `json:"sender_id"`  // ID of the sending client
}

func NewMessage(msgType MessageType, payload map[string]any) *Message {
    return &Message{
        Type:      msgType,
        Timestamp: time.Now(),
        Payload:   payload,
    }
}

func (m *Message) Encode() ([]byte, error) {
    return json.Marshal(m)
}

func DecodeMessage(data []byte) (*Message, error) {
    var msg Message
    if err := json.Unmarshal(data, &msg); err != nil {
        return nil, err
    }
    return &msg, nil
}