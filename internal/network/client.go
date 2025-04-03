package network

import (
    "net"
    "time"
    "sync"
    "os"
    "github.com/shivangjhalani/remotesync/internal/logger"
    "github.com/shivangjhalani/remotesync/internal/protocol"
)

type ClientConfig struct {
    ServerHost string
    ServerPort string
    RetryDelay time.Duration
    MaxRetries int
    ClientName string
    ClientID   string
    MaxBandwidth int64
    MonitorID    int
}

type RemoteClient struct {
    conn       net.Conn
    config     ClientConfig
    connected  bool
    mutex      sync.RWMutex
    retryCount int
    clientName string
    clientID   string
    mouseCapture *input.MouseCapture
    mouseSim     *input.MouseSimulator
    keyboardCapture *input.KeyboardCapture
    keyboardSim    *input.KeyboardSimulator
    screenTransmitter *screen.ScreenTransmitter
    tokenManager *security.TokenManager
    encryptor    *security.Encryptor
}

func NewRemoteClient(config ClientConfig) *RemoteClient {
    if config.RetryDelay == 0 {
        config.RetryDelay = 5 * time.Second
    }
    if config.MaxRetries == 0 {
        config.MaxRetries = 5
    }
    
    return &RemoteClient{
        config: config,
    }
}

func (c *RemoteClient) Connect() error {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    if c.connected {
        return nil
    }

    addr := net.JoinHostPort(c.config.ServerHost, c.config.ServerPort)
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        return err
    }

    c.conn = conn
    c.connected = true
    c.retryCount = 0
    
    if err := c.identifyClient(); err != nil {
        c.conn.Close()
        c.connected = false
        return err
    }
    
    // Authenticate
    if err := c.authenticate(); err != nil {
        c.conn.Close()
        return fmt.Errorf("authentication failed: %v", err)
    }
    
    go c.startHeartbeat()
    
    return nil
}

func (c *RemoteClient) identifyClient() error {
    hostname, _ := os.Hostname()
    if c.config.ClientName == "" {
        c.config.ClientName = hostname
    }

    payload := map[string]any{
        "name": c.config.ClientName,
        "hostname": hostname,
    }

    msg := protocol.NewMessage(protocol.TypeIdentify, payload)
    data, err := msg.Encode()
    if err != nil {
        return err
    }

    _, err = c.conn.Write(data)
    return err
}

func (c *RemoteClient) authenticate() error {
    creds := security.Credentials{
        Username: c.config.Username,
        Password: c.config.Password,
    }

    msg := protocol.NewMessage(protocol.TypeAuth, map[string]any{
        "username": creds.Username,
        "password": security.HashPassword(creds.Password),
    })

    if err := c.sendMessage(msg); err != nil {
        return err
    }

    // Wait for auth response
    // TODO: Implement response handling

    return nil
}

func (c *RemoteClient) startHeartbeat() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        <-ticker.C
        if err := c.sendPing(); err != nil {
            logger.ErrorLogger.Printf("Failed to send heartbeat: %v", err)
            c.reconnect()
        }
    }
}

func (c *RemoteClient) sendPing() error {
    c.mutex.RLock()
    defer c.mutex.RUnlock()

    if !c.connected {
        return fmt.Errorf("client not connected")
    }

    payload := map[string]any{
        "timestamp": time.Now(),
    }

    msg := protocol.NewMessage(protocol.TypePing, payload)
    msg.SenderID = c.clientID
    
    data, err := msg.Encode()
    if err != nil {
        return err
    }

    _, err = c.conn.Write(data)
    return err
}

func (c *RemoteClient) reconnect() {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    if c.retryCount >= c.config.MaxRetries {
        logger.ErrorLogger.Fatal("Max retry attempts reached")
        return
    }

    c.connected = false
    if c.conn != nil {
        c.conn.Close()
    }

    time.Sleep(c.config.RetryDelay)
    c.retryCount++

    logger.InfoLogger.Printf("Attempting reconnection (attempt %d/%d)", 
        c.retryCount, c.config.MaxRetries)

    if err := c.Connect(); err != nil {
        logger.ErrorLogger.Printf("Reconnection failed: %v", err)
        go c.reconnect()
    }
}

func (c *RemoteClient) StartInputCapture() {
    c.mouseCapture = input.NewMouseCapture()
    c.mouseSim = input.NewMouseSimulator()
    
    c.mouseCapture.Start()
    go c.handleMouseEvents()
    
    c.keyboardCapture = input.NewKeyboardCapture()
    c.keyboardSim = input.NewKeyboardSimulator()
    
    c.keyboardCapture.Start()
    go c.handleKeyboardEvents()
}

func (c *RemoteClient) handleMouseEvents() {
    for event := range c.mouseCapture.eventChan {
        payload := map[string]any{
            "x":      event.X,
            "y":      event.Y,
            "button": event.Button,
            "action": event.Action,
        }

        msg := protocol.NewMessage(protocol.TypeMouseEvent, payload)
        msg.SenderID = c.clientID

        if err := c.sendMessage(msg); err != nil {
            logger.ErrorLogger.Printf("Failed to send mouse event: %v", err)
        }
    }
}

func (c *RemoteClient) handleKeyboardEvents() {
    for event := range c.keyboardCapture.eventChan {
        payload := map[string]any{
            "key_code": event.KeyCode,
            "key_char": event.KeyChar,
            "modifiers": event.Modifiers,
            "action": event.Action,
        }

        msg := protocol.NewMessage(protocol.TypeKeyEvent, payload)
        msg.SenderID = c.clientID

        if err := c.sendMessage(msg); err != nil {
            logger.ErrorLogger.Printf("Failed to send keyboard event: %v", err)
        }
    }
}

func (c *RemoteClient) StartScreenSharing(config screen.CaptureConfig) error {
    config.MaxBandwidth = c.config.MaxBandwidth
    config.MonitorID = c.config.MonitorID
    
    c.screenTransmitter = screen.NewScreenTransmitter(config, c.sendMessage)
    return c.screenTransmitter.Start()
}

func (c *RemoteClient) StopScreenSharing() {
    if c.screenTransmitter != nil {
        c.screenTransmitter.Stop()
    }
}

func (c *RemoteClient) sendMessage(msg *protocol.Message) error {
    data, err := msg.Encode()
    if (err != nil) {
        return err
    }

    // Encrypt data before sending
    encrypted, err := c.encryptor.Encrypt(data)
    if (err != nil) {
        return err
    }

    _, err = c.conn.Write(encrypted)
    return err
}