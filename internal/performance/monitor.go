package performance

import (
    "sync"
    "time"
)

type Metric struct {
    Name      string
    Value     float64
    Timestamp time.Time
}

type PerformanceMonitor struct {
    metrics   map[string][]Metric
    mutex     sync.RWMutex
}

func NewPerformanceMonitor() *PerformanceMonitor {
    return &PerformanceMonitor{
        metrics: make(map[string][]Metric),
    }
}

func (pm *PerformanceMonitor) RecordMetric(name string, value float64) {
    pm.mutex.Lock()
    defer pm.mutex.Unlock()

    metric := Metric{
        Name:      name,
        Value:     value,
        Timestamp: time.Now(),
    }

    pm.metrics[name] = append(pm.metrics[name], metric)
}

func (pm *PerformanceMonitor) GetAverageMetric(name string) float64 {
    pm.mutex.RLock()
    defer pm.mutex.RUnlock()

    metrics := pm.metrics[name]
    if len(metrics) == 0 {
        return 0
    }

    var sum float64
    for _, m := range metrics {
        sum += m.Value
    }
    return sum / float64(len(metrics))
}