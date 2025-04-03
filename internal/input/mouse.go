package input

import (
    "github.com/go-vgo/robotgo"
    "github.com/shivangjhalani/remotesync/internal/protocol"
    "github.com/shivangjhalani/remotesync/internal/logger"
    "image"
)

type MouseEvent struct {
    X      int     `json:"x"`
    Y      int     `json:"y"`
    Button string  `json:"button"` // left, right, middle
    Action string  `json:"action"` // move, click, scroll
    Delta  int     `json:"delta"`  // for scroll events
}

type MouseCapture struct {
    enabled     bool
    eventChan   chan *MouseEvent
}

func NewMouseCapture() *MouseCapture {
    return &MouseCapture{
        eventChan: make(chan *MouseEvent, 100),
    }
}

func (m *MouseCapture) Start() {
    m.enabled = true
    go m.captureMouseEvents()
}

func (m *MouseCapture) Stop() {
    m.enabled = false
    close(m.eventChan)
}

func (m *MouseCapture) captureMouseEvents() {
    for m.enabled {
        x, y := robotgo.GetMousePos()
        
        // Check for mouse movement
        event := &MouseEvent{
            X: x,
            Y: y,
            Action: "move",
        }
        
        // Check for mouse clicks
        if robotgo.GetMouseButton("left") {
            event.Button = "left"
            event.Action = "click"
        } else if robotgo.GetMouseButton("right") {
            event.Button = "right"
            event.Action = "click"
        }
        
        m.eventChan <- event
    }
}

type MouseSimulator struct {
    cursorManager *CursorManager
}

func (ms *MouseSimulator) SimulateEvent(event *MouseEvent, clientID string) error {
    cursor := ms.cursorManager.cursors[clientID]
    if cursor == nil {
        cursor = ms.cursorManager.AddCursor(clientID, "default")
    }

    cursor.Position = image.Point{X: event.X, Y: event.Y}
    
    // Only simulate actual mouse events for active controller
    if ms.isActiveController(clientID) {
        switch event.Action {
        case "move":
            robotgo.MoveMouse(event.X, event.Y)
        case "click":
            robotgo.MouseClick(event.Button)
        case "scroll":
            robotgo.ScrollMouse(event.Delta, "vertical")
        }
    }

    return nil
}