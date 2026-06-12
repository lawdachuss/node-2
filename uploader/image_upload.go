package uploader

import (
	"fmt"
	"sync"
	"time"
)

// pixhostSem limits concurrent Pixhost.to uploads to avoid rate limiting.
var pixhostSem = make(chan struct{}, 5)

// MultiImageUploader uploads thumbnails/sprites to all configured hosts
// in parallel. Prefers Pixhost.to, falls back to ImgBB.
type MultiImageUploader struct {
	pixhost *ThumbnailUploader
	imgbb   *ImgBBUploader
}

// NewMultiImageUploader creates a new image uploader that uploads to
// Pixhost.to and ImgBB simultaneously.
func NewMultiImageUploader() *MultiImageUploader {
	return &MultiImageUploader{
		pixhost: NewThumbnailUploader(""),
		imgbb:   NewImgBBUploader(),
	}
}

// Upload uploads to Pixhost (with retries) and ImgBB in parallel.
// Returns Pixhost URL on success, ImgBB URL as fallback.
func (m *MultiImageUploader) Upload(filePath string) (url, host string, err error) {
	var (
		mu         sync.Mutex
		pixhostURL string
		pixhostErr error
		imgbbURL   string
		imgbbErr   error
		wg         sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		pixhostSem <- struct{}{}
		defer func() { <-pixhostSem }()
		var lastErr error
		for attempt := 0; attempt < 3; attempt++ {
			if attempt > 0 {
				time.Sleep(time.Duration(1<<attempt) * time.Second)
			}
			u, err := m.pixhost.Upload(filePath)
			if err == nil {
				mu.Lock()
				pixhostURL = u
				mu.Unlock()
				return
			}
			lastErr = err
		}
		mu.Lock()
		pixhostErr = lastErr
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		u, err := m.imgbb.Upload(filePath)
		mu.Lock()
		imgbbURL = u
		imgbbErr = err
		mu.Unlock()
	}()

	wg.Wait()

	if pixhostURL != "" {
		return pixhostURL, "Pixhost", nil
	}
	if imgbbURL != "" {
		return imgbbURL, "ImgBB", nil
	}
	return "", "", fmt.Errorf("pixhost: %w (imgbb also failed: %v)", pixhostErr, imgbbErr)
}
