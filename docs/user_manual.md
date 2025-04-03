# RemoteSync User Manual

## Getting Started

### Installation
1. Download the latest release package
2. Extract the archive
3. Run the server or client executable

### Server Setup
1. Start the server:
```bash
./remotesync-server -port 8080
```

2. Monitor connected clients through the server UI

### Client Usage
1. Start the client:
```bash
./remotesync-client -host localhost -port 8080
```

2. Connect to available clients:
   - Select a client from the available list
   - Choose control mode (Exclusive/Simultaneous)
   - Click "Request Control"

### Control Modes
- **Exclusive Mode**: Only one client can control at a time
- **Simultaneous Mode**: Multiple clients can control simultaneously

### Settings
- Frame Rate: Adjust screen capture frame rate
- Quality: Adjust screen capture quality
- Monitor: Select which monitor to share