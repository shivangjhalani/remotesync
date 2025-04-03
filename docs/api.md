# RemoteSync API Documentation

## Network Protocol

### Message Structure
```json
{
    "type": "message_type",
    "timestamp": "2024-04-03T12:00:00Z",
    "payload": {},
    "sender_id": "client_id"
}
```

### Message Types
- `identify`: Client identification
- `auth`: Authentication request
- `control_request`: Request control of another client
- `mouse_event`: Mouse input events
- `key_event`: Keyboard input events
- `screen_data`: Screen capture data

## Client API

### RemoteClient
```go
type ClientConfig struct {
    ServerHost    string
    ServerPort    string
    RetryDelay    time.Duration
    MaxRetries    int
    MaxBandwidth  int64
    MonitorID     int
}

func NewRemoteClient(config ClientConfig) *RemoteClient
func (c *RemoteClient) Connect() error
func (c *RemoteClient) StartScreenSharing(config screen.CaptureConfig) error
func (c *RemoteClient) RequestControl(targetID string, mode ControlMode) error
```

## Server API

### Server
```go
func NewServer() *Server
func (s *Server) Start(port string) error
func (s *Server) RegisterClient(conn net.Conn) (*Client, error)
func (s *Server) GetClientList() []*Client
```