package session

type ControlRequest struct {
    RequesterID  string
    TargetID     string
    Mode         ControlMode
}

type ControlResponse struct {
    Approved     bool
    SessionID    string
    Message      string
}

func (sm *SessionManager) RequestControl(req *ControlRequest) *ControlResponse {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    // Check if target is already being controlled
    for _, session := range sm.sessions {
        if session.ControlledID == req.TargetID && session.Active {
            if session.Mode == ModeExclusive {
                return &ControlResponse{
                    Approved: false,
                    Message: "Target is already being controlled exclusively",
                }
            }
            // Add controller to existing simultaneous session
            if session.Mode == ModeSimultaneous {
                session.Controllers[req.RequesterID] = true
                return &ControlResponse{
                    Approved:  true,
                    SessionID: session.ID,
                    Message:   "Added to existing simultaneous control session",
                }
            }
        }
    }

    // Create new session
    sessionID := generateSessionID()
    session := &Session{
        ID:           sessionID,
        ControlledID: req.TargetID,
        Controllers:  map[string]bool{req.RequesterID: true},
        Mode:        req.Mode,
        StartTime:   time.Now(),
        Active:      true,
    }

    sm.sessions[sessionID] = session
    return &ControlResponse{
        Approved:  true,
        SessionID: sessionID,
        Message:   "New control session created",
    }
}