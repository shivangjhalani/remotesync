package screen

import (
    "time"
    "sync"
)

type BandwidthManager struct {
    bytesTransferred int64
    lastReset       time.Time
    mutex           sync.Mutex
    maxBandwidth    int64
}

func NewBandwidthManager(maxBandwidth int64) *BandwidthManager {
    return &BandwidthManager{
        lastReset:    time.Now(),
        maxBandwidth: maxBandwidth,
    }
}

func (bm *BandwidthManager) TrackTransfer(bytes int64) bool {
    bm.mutex.Lock()
    defer bm.mutex.Unlock()

    now := time.Now()
    if now.Sub(bm.lastReset) > time.Second {
        bm.bytesTransferred = 0
        bm.lastReset = now
    }

    if bm.bytesTransferred+bytes > bm.maxBandwidth {
        return false
    }

    bm.bytesTransferred += bytes
    return true
}