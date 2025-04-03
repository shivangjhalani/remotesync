package screen

import (
    "remotesync/internal/protocol"
    "remotesync/internal/logger"
)

type ScreenTransmitter struct {
    capture    *ScreenCapture
    sendFrame  func(*protocol.Message) error
}

func NewScreenTransmitter(frameRate, quality int, sender func(*protocol.Message) error) *ScreenTransmitter {
    return &ScreenTransmitter{
        capture:   NewScreenCapture(frameRate, quality),
        sendFrame: sender,
    }
}

func (st *ScreenTransmitter) Start() error {
    if err := st.capture.Start(); err != nil {
        return err
    }

    go st.transmitFrames()
    return nil
}

func (st *ScreenTransmitter) transmitFrames() {
    for frame := range st.capture.frameChan {
        payload := map[string]any{
            "image":      frame.Image,
            "width":      frame.Resolution.X,
            "height":     frame.Resolution.Y,
            "timestamp":  frame.Timestamp,
        }

        msg := protocol.NewMessage(protocol.TypeScreenData, payload)
        if err := st.sendFrame(msg); err != nil {
            logger.ErrorLogger.Printf("Failed to send frame: %v", err)
        }
    }
}