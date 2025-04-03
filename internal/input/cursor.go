package input

import (
    "image"
    "image/color"
)

type CursorStyle struct {
    Color     color.RGBA
    Size      int
    Label     string
}

type CursorManager struct {
    cursors    map[string]*Cursor
    styles     map[string]CursorStyle
}

func NewCursorManager() *CursorManager {
    return &CursorManager{
        cursors: make(map[string]*Cursor),
        styles: map[string]CursorStyle{
            "default": {
                Color: color.RGBA{R: 255, G: 0, B: 0, A: 255},
                Size:  20,
            },
            "secondary": {
                Color: color.RGBA{R: 0, G: 0, B: 255, A: 255},
                Size:  20,
            },
        },
    }
}

func (cm *CursorManager) AddCursor(clientID string, style string) *Cursor {
    cursorStyle := cm.styles[style]
    if cursorStyle == (CursorStyle{}) {
        cursorStyle = cm.styles["default"]
    }

    cursor := &Cursor{
        ID:       clientID,
        Style:    cursorStyle,
        Position: image.Point{X: 0, Y: 0},
        Visible:  true,
    }

    cm.cursors[clientID] = cursor
    return cursor
}

type CursorVisualizer struct {
    cursors map[string]*Cursor
    overlay *screen.Overlay
}

func NewCursorVisualizer() *CursorVisualizer {
    return &CursorVisualizer{
        cursors: make(map[string]*Cursor),
        overlay: screen.NewOverlay(),
    }
}

func (cv *CursorVisualizer) DrawCursors() {
    cv.overlay.Clear()
    for _, cursor := range cv.cursors {
        if cursor.Visible {
            cv.overlay.DrawCursor(cursor.Position, cursor.Style.Color, cursor.Style.Size)
            cv.overlay.DrawLabel(cursor.Position, cursor.ID, cursor.Style.Color)
        }
    }
    cv.overlay.Update()
}