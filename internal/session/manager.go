package session

import (
    "sync"
    "time"
    "remotesync/internal/logger"
    "fmt"
    "remotesync/internal/input"
    "image"
)

type ControlMode string

const (
    ModeExclusive    ControlMode = "exclusive"
    ModeSimultaneous ControlMode = "simultaneous"
)

type Session struct {
    ID                string
    ControlledID      string
    Controllers       map[string]bool
    Mode              ControlMode
    StartTime         time.Time
    Active            bool
    mutex             sync.RWMutex
    CursorManager     *input.CursorManager
    CursorVisualizer  *input.CursorVisualizer
    KeyboardManager   *input.MultiKeyboardManager
}

type SessionManager struct {
    sessions     map[string]*Session
    mutex        sync.RWMutex
}

func NewSessionManager() *SessionManager {
    return &SessionManager{
        sessions: make(map[string]*Session),
    }
}

func (sm *SessionManager) TerminateSession(sessionID string, requesterID string) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    session, exists := sm.sessions[sessionID]
    if (!exists) {
        return fmt.Errorf("session not found")
    }

    // Check if requester has permission to terminate
    if (!session.Controllers[requesterID] && session.ControlledID != requesterID) {
        return fmt.Errorf("no permission to terminate session")
    }

    session.Active = false
    delete(sm.sessions, sessionID)
    return nil
}

func (sm *SessionManager) RequestControl(req *ControlRequest) *ControlResponse {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    session := &Session{
        ID:                req.SessionID,
        ControlledID:      req.ControlledID,
        Controllers:       make(map[string]bool),
        Mode:              req.Mode,
        StartTime:         time.Now(),
        Active:            true,
        CursorManager:     input.NewCursorManager(),
        CursorVisualizer:  input.NewCursorVisualizer(),
        KeyboardManager:   input.NewMultiKeyboardManager(),
    }

    // Assign different cursor styles to controllers
    if req.Mode == ModeSimultaneous {
        session.CursorManager.AddCursor(req.RequesterID, "secondary")
    } else {
        session.CursorManager.AddCursor(req.RequesterID, "default")
    }

    sm.sessions[req.SessionID] = session

    return &ControlResponse{
        SessionID: req.SessionID,
        Success:   true,
    }
}

func (sm *SessionManager) HandleInput(sessionID string, clientID string, input interface{}) error {
    session, exists := sm.sessions[sessionID]
    if !exists {
        return fmt.Errorf("session not found")
    }

    switch evt := input.(type) {
    case *input.MouseEvent:
        cursor := session.CursorVisualizer.GetCursor(clientID)
        if cursor == nil {
            return fmt.Errorf("cursor not found")
        }
        cursor.Position = image.Point{X: evt.X, Y: evt.Y}
        session.CursorVisualizer.DrawCursors()

    case *input.KeyEvent:
        return session.KeyboardManager.HandleKeyEvent(evt, clientID)
    }

    return nil
}