package uploader

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// ThumbnailUploader handles uploading thumbnail images to external hosting
type ThumbnailUploader struct {
	client *http.Client
}

// NewThumbnailUploader creates a new thumbnail uploader using Catbox.moe
func NewThumbnailUploader() *ThumbnailUploader {
	return &ThumbnailUploader{
		client: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

// Upload uploads a thumbnail image to Catbox.moe and returns the direct image URL
// Catbox.moe allows NSFW content and provides permanent hosting
func (t *ThumbnailUploader) Upload(thumbnailPath string) (string, error) {
	log.Printf("Uploading thumbnail to Catbox.moe: %s", thumbnailPath)
	
	// Open the image file
	file, err := os.Open(thumbnailPath)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()
	
	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// Add reqtype field
	if err := writer.WriteField("reqtype", "fileupload"); err != nil {
		return "", fmt.Errorf("write reqtype field: %w", err)
	}
	
	// Add image file
	part, err := writer.CreateFormFile("fileToUpload", filepath.Base(thumbnailPath))
	if err != nil {
		return "", fmt.Errorf("create form file: %w", err)
	}
	
	// Copy file content
	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("copy file: %w", err)
	}
	
	// Close writer to finalize multipart form
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("close writer: %w", err)
	}
	
	// Create request to Catbox API
	uploadURL := "https://catbox.moe/user/api.php"
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	// Send request
	resp, err := t.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(respBody))
	}
	
	// Catbox returns the direct URL as plain text
	imageURL := string(respBody)
	if imageURL == "" {
		return "", fmt.Errorf("empty response from Catbox")
	}
	
	log.Printf("Thumbnail uploaded successfully: %s", imageURL)
	return imageURL, nil
}
