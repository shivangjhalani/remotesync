package input

import (
    "github.com/go-vgo/robotgo"
    "github.com/shivangjhalani/remotesync/internal/logger"
)

type MouseSimulator struct {
    cursors map[string]*Cursor
}

type Cursor struct {
    ID     string
    X      int
    Y      int
    Visible bool
}

func NewMouseSimulator() *MouseSimulator {
    return &MouseSimulator{
        cursors: make(map[string]*Cursor),
    }
}

func (ms *MouseSimulator) SimulateEvent(event *MouseEvent, cursorID string) error {
    cursor, exists := ms.cursors[cursorID]
    if !exists {
        cursor = &Cursor{ID: cursorID}
        ms.cursors[cursorID] = cursor
    }

    cursor.X = event.X
    cursor.Y = event.Y

    switch event.Action {
    case "move":
        robotgo.MoveMouse(event.X, event.Y)
    case "click":
        robotgo.MouseClick(event.Button)
    case "scroll":
        robotgo.ScrollMouse(event.Delta, "vertical")
    }

    return nil
}