package security_test

import (
    "testing"
    "github.com/shivangjhalani/remotesync/internal/security"
)

func TestTokenGeneration(t *testing.T) {
    tm := security.NewTokenManager("test-secret")
    
    token, err := tm.GenerateToken("testuser")
    if err != nil {
        t.Errorf("Failed to generate token: %v", err)
    }

    username, err := tm.ValidateToken(token)
    if err != nil {
        t.Errorf("Failed to validate token: %v", err)
    }

    if username != "testuser" {
        t.Errorf("Username mismatch: got %v, want %v", 
                 username, "testuser")
    }
}