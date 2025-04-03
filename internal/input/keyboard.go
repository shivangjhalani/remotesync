package input

import (
    "github.com/go-vgo/robotgo"
    "remotesync/internal/logger"
)

type KeyEvent struct {
    KeyCode    string `json:"key_code"`
    KeyChar    string `json:"key_char"`
    Modifiers  []string `json:"modifiers"` // ctrl, alt, shift, etc.
    Action     string `json:"action"` // press, release
}

type KeyboardCapture struct {
    enabled   bool
    eventChan chan *KeyEvent
}

func NewKeyboardCapture() *KeyboardCapture {
    return &KeyboardCapture{
        eventChan: make(chan *KeyEvent, 100),
    }
}

func (k *KeyboardCapture) Start() {
    k.enabled = true
    go k.captureKeyEvents()
}

func (k *KeyboardCapture) Stop() {
    k.enabled = false
    close(k.eventChan)
}

func (k *KeyboardCapture) captureKeyEvents() {
    for k.enabled {
        robotgo.AddEvents("k")  // Listen for keyboard events
        if key := robotgo.GetKey(); key != "" {
            event := &KeyEvent{
                KeyCode: key,
                Action: "press",
            }
            
            // Check modifiers
            if robotgo.GetKey("ctrl") {
                event.Modifiers = append(event.Modifiers, "ctrl")
            }
            if robotgo.GetKey("alt") {
                event.Modifiers = append(event.Modifiers, "alt")
            }
            if robotgo.GetKey("shift") {
                event.Modifiers = append(event.Modifiers, "shift")
            }
            
            k.eventChan <- event
        }
    }
}