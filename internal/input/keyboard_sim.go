package input

import (
    "github.com/go-vgo/robotgo"
    "remotesync/internal/logger"
    "sync"
)

type KeyboardSimulator struct {
    activeModifiers map[string]bool
}

func NewKeyboardSimulator() *KeyboardSimulator {
    return &KeyboardSimulator{
        activeModifiers: make(map[string]bool),
    }
}

func (ks *KeyboardSimulator) SimulateEvent(event *KeyEvent) error {
    // Handle modifiers first
    for _, mod := range event.Modifiers {
        if !ks.activeModifiers[mod] {
            robotgo.KeyToggle(mod, "down")
            ks.activeModifiers[mod] = true
        }
    }

    // Simulate the key press/release
    if event.Action == "press" {
        robotgo.KeyTap(event.KeyCode)
    } else {
        robotgo.KeyToggle(event.KeyCode, "up")
    }

    // Release modifiers
    for mod := range ks.activeModifiers {
        robotgo.KeyToggle(mod, "up")
        delete(ks.activeModifiers, mod)
    }

    return nil
}

type MultiKeyboardManager struct {
    simulators map[string]*KeyboardSimulator
    mutex      sync.RWMutex
}

func NewMultiKeyboardManager() *MultiKeyboardManager {
    return & MultiKeyboardManager{
        simulators: make(map[string]*KeyboardSimulator),
    }
}

func (mkm *MultiKeyboardManager) HandleKeyEvent(event *KeyEvent, clientID string) error {
    mkm.mutex.Lock()
    defer mkm.mutex.Unlock()

    sim, exists := mkm.simulators[clientID]
    if !exists {
        sim = NewKeyboardSimulator()
        mkm.simulators[clientID] = sim
    }

    return sim.SimulateEvent(event)
}