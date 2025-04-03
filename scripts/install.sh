#!/bin/bash

# Installation script for RemoteSync

# Check for required dependencies
check_dependencies() {
    echo "Checking dependencies..."
    command -v go >/dev/null 2>&1 || { echo "Go is required but not installed. Aborting." >&2; exit 1; }
    command -v make >/dev/null 2>&1 || { echo "Make is required but not installed. Aborting." >&2; exit 1; }
}

# Install system dependencies
install_dependencies() {
    echo "Installing system dependencies..."
    if [ -f /etc/debian_version ]; then
        sudo apt-get update
        sudo apt-get install -y libopencv-dev libx11-dev
    elif [ -f /etc/redhat-release ]; then
        sudo dnf install -y opencv-devel libX11-devel
    fi
}

# Build and install RemoteSync
install_remotesync() {
    echo "Building RemoteSync..."
    make clean
    make build
    
    echo "Installing RemoteSync..."
    sudo make install
}

# Create configuration directory
setup_config() {
    echo "Setting up configuration..."
    sudo mkdir -p /etc/remotesync
    sudo cp config/*.json /etc/remotesync/
}

# Main installation process
main() {
    check_dependencies
    install_dependencies
    install_remotesync
    setup_config
    
    echo "Installation completed successfully!"
    echo "Run 'remotesync-server' to start the server"
    echo "Run 'remotesync-client' to start the client"
}

main "$@"