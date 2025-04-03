#!/bin/bash

# System test script for RemoteSync

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# Test server startup
test_server() {
    echo "Testing server startup..."
    ./bin/remotesync-server &
    SERVER_PID=$!
    sleep 2
    
    if ps -p $SERVER_PID > /dev/null; then
        echo -e "${GREEN}Server startup test: PASSED${NC}"
        return 0
    else
        echo -e "${RED}Server startup test: FAILED${NC}"
        return 1
    fi
}

# Test client connection
test_client_connection() {
    echo "Testing client connection..."
    ./bin/remotesync-client -host localhost -port 8080 &
    CLIENT_PID=$!
    sleep 2
    
    if ps -p $CLIENT_PID > /dev/null; then
        echo -e "${GREEN}Client connection test: PASSED${NC}"
        kill $CLIENT_PID
        return 0
    else
        echo -e "${RED}Client connection test: FAILED${NC}"
        return 1
    fi
}

# Run all tests
run_tests() {
    make clean
    make build
    
    test_server
    SERVER_RESULT=$?
    
    test_client_connection
    CLIENT_RESULT=$?
    
    # Cleanup
    kill $SERVER_PID
    
    # Final results
    echo "Test Results:"
    echo "-------------"
    [ $SERVER_RESULT -eq 0 ] && echo -e "Server Tests: ${GREEN}PASSED${NC}" || echo -e "Server Tests: ${RED}FAILED${NC}"
    [ $CLIENT_RESULT -eq 0 ] && echo -e "Client Tests: ${GREEN}PASSED${NC}" || echo -e "Client Tests: ${RED}FAILED${NC}"
}

run_tests