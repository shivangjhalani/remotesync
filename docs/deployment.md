# RemoteSync Deployment Guide

## System Requirements
- Go 1.19 or later
- OpenCV for screen capture
- X11 development libraries

## Building from Source
1. Clone the repository:
```bash
git clone https://github.com/shivangjhalani/remotesync.git
```

2. Install dependencies:
```bash
sudo apt-get install libopencv-dev
go mod download
```

3. Build the executables:
```bash
make build
```

## Deployment
1. Copy the built executables to the target system
2. Set up configuration files
3. Start the server service
4. Distribute client executables

## Security Considerations
- Configure firewall rules
- Set up authentication
- Use TLS for secure communication