package screen

import (
    "gocv.io/x/gocv"
    "image"
    "remotesync/internal/logger"
)

type QualityLevel int

const (
    QualityLow    QualityLevel = 50
    QualityMedium QualityLevel = 75
    QualityHigh   QualityLevel = 90
)

type CaptureConfig struct {
    FrameRate     int
    Quality       QualityLevel
    MonitorID     int
    MaxBandwidth  int64 // bytes per second
}

type ScreenCapture struct {
    enabled      bool
    frameChan    chan *Frame
    window       *gocv.Window
    capture      *gocv.VideoCapture
    config       CaptureConfig
    monitorMgr   *MonitorManager
    bandwidthMgr *BandwidthManager
}

type Frame struct {
    Image       []byte
    Resolution  image.Point
    Timestamp   int64
}

func NewScreenCapture(frameRate, quality int) *ScreenCapture {
    if frameRate == 0 {
        frameRate = 30
    }
    if quality == 0 {
        quality = 75
    }

    return &ScreenCapture{
        frameChan: make(chan *Frame, 10),
        config: CaptureConfig{
            FrameRate: frameRate,
            Quality:   QualityLevel(quality),
        },
    }
}

func (sc *ScreenCapture) Start() error {
    sc.enabled = true
    go sc.captureScreen()
    return nil
}

func (sc *ScreenCapture) Stop() {
    sc.enabled = false
    close(sc.frameChan)
    if sc.capture != nil {
        sc.capture.Close()
    }
}

func (sc *ScreenCapture) captureScreen() {
    for sc.enabled {
        img := gocv.NewMat()
        defer img.Close()
        
        // Capture screen using gocv
        screen := gocv.IMRead("screen", gocv.IMReadColor)
        if screen.Empty() {
            logger.ErrorLogger.Println("Failed to capture screen")
            continue
        }

        // Compress image
        buf, err := gocv.IMEncode(".jpg", screen, []int{gocv.IMWriteJpegQuality, int(sc.config.Quality)})
        if err != nil {
            logger.ErrorLogger.Printf("Failed to compress frame: %v", err)
            continue
        }

        frame := &Frame{
            Image:      buf,
            Resolution: image.Point{X: screen.Cols(), Y: screen.Rows()},
            Timestamp:  time.Now().UnixNano(),
        }

        sc.frameChan <- frame
        sc.adjustQuality(int64(len(buf)))
        time.Sleep(time.Second / time.Duration(sc.config.FrameRate))
    }
}

func (sc *ScreenCapture) adjustQuality(frameSize int64) {
    // Dynamic quality adjustment based on bandwidth usage
    if frameSize > sc.config.MaxBandwidth/int64(sc.config.FrameRate) {
        if sc.config.Quality > QualityLow {
            sc.config.Quality -= 10
        }
        if sc.config.FrameRate > 15 {
            sc.config.FrameRate--
        }
    } else if frameSize < sc.config.MaxBandwidth/(2*int64(sc.config.FrameRate)) {
        if sc.config.Quality < QualityHigh {
            sc.config.Quality += 5
        }
        if sc.config.FrameRate < 30 {
            sc.config.FrameRate++
        }
    }
}