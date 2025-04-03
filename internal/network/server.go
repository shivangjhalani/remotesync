package network

import (
    "bufio"
    "fmt"
    "net"
    "sync"
    "time"
    "crypto/rand"
    "encoding/hex"
    "github.com/shivangjhalani/remotesync/internal/logger"
    "github.com/shivangjhalani/remotesync/internal/protocol"
    "github.com/shivangjhalani/remotesync/internal/session"
    "github.com/shivangjhalani/remotesync/internal/security"
    "github.com/shivangjhalani/remotesync/internal/performance"
)

type Server struct {
    clients        map[string]*Client
    mutex          sync.RWMutex
    maxClients     int
    sessionManager *session.SessionManager
    tokenManager   *security.TokenManager
    encryptor      *security.Encryptor
    perfMonitor    *performance.PerformanceMonitor
}

type Client struct {
    ID        string
    Name      string
    Conn      net.Conn
    LastPing  time.Time
}

func NewServer() *Server {
    return &Server{
        clients:    make(map[string]*Client),
        maxClients: 2, // As per requirement of 2 clients
    }
}

func (s *Server) RegisterClient(conn net.Conn) (*Client, error) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    if len(s.clients) >= s.maxClients {
        return nil, fmt.Errorf("maximum number of clients reached")
    }

    client := &Client{
        ID:       generateID(),
        Conn:     conn,
        LastPing: time.Now(),
    }

    s.clients[client.ID] = client
    return client, nil
}

func generateID() string {
    bytes := make([]byte, 8)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}

func (s *Server) GetClientList() []*Client {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    
    clients := make([]*Client, 0, len(s.clients))
    for _, client := range s.clients {
        clients = append(clients, client)
    }
    return clients
}

func (s *Server) RemoveClient(id string) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    if client, exists := s.clients[id]; exists {
        client.Conn.Close()
        delete(s.clients, id)
        logger.InfoLogger.Printf("Client %s disconnected", id)
    }
}

func (s *Server) BroadcastClientList() {
    // TODO: Implement when communication protocol is ready
}

func (s *Server) handleIdentification(conn net.Conn) (*Client, error) {
    reader := bufio.NewReader(conn)
    data, err := reader.ReadBytes('\n')
    if err != nil {
        return nil, err
    }

    msg, err := protocol.DecodeMessage(data)
    if err != nil || msg.Type != protocol.TypeIdentify {
        return nil, fmt.Errorf("invalid identification message")
    }

    name, ok := msg.Payload["name"].(string)
    if !ok {
        return nil, fmt.Errorf("invalid client name")
    }

    client := &Client{
        ID:        generateID(),
        Name:      name,
        Conn:      conn,
        LastPing:  time.Now(),
    }

    return client, nil
}

func (s *Server) handleControlRequest(msg *protocol.Message) error {
    req := &session.ControlRequest{
        RequesterID: msg.SenderID,
        TargetID:    msg.Payload["target_id"].(string),
        Mode:        session.ControlMode(msg.Payload["mode"].(string)),
    }

    response := s.sessionManager.RequestControl(req)
    
    // Send response back to requester
    responseMsg := protocol.NewMessage(protocol.TypeControlRes, map[string]any{
        "approved":   response.Approved,
        "session_id": response.SessionID,
        "message":    response.Message,
    })
    
    return s.sendToClient(msg.SenderID, responseMsg)
}

func (s *Server) handleAuth(msg *protocol.Message) error {
    username := msg.Payload["username"].(string)
    password := msg.Payload["password"].(string)

    // In a real application, validate against a database
    // For now, use a simple check
    if !s.validateCredentials(username, password) {
        return fmt.Errorf("invalid credentials")
    }

    token, err := s.tokenManager.GenerateToken(username)
    if err != nil {
        return err
    }

    response := protocol.NewMessage(protocol.TypeAuthResponse, map[string]any{
        "token": token,
        "success": true,
    })

    return s.sendToClient(msg.SenderID, response)
}

func (s *Server) handleMessage(msg *protocol.Message) error {
    start := time.Now()
    err := s.processMessage(msg)
    duration := time.Since(start)
    
    s.perfMonitor.RecordMetric("message_processing_time", float64(duration.Milliseconds()))
    return err
}