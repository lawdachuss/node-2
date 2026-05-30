package uploader

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// StreamtapeUploader handles uploading files to Streamtape
type StreamtapeUploader struct {
	login  string
	key    string
	client *http.Client
}

// NewStreamtapeUploader creates a new Streamtape uploader instance
func NewStreamtapeUploader(login, key string) *StreamtapeUploader {
	return &StreamtapeUploader{
		login: login,
		key:   key,
		client: &http.Client{
			Timeout: 120 * time.Minute,
			Transport: &http.Transport{
				MaxIdleConns:          10,
				MaxIdleConnsPerHost:   2,
				IdleConnTimeout:       90 * time.Second,
				DisableCompression:    true,
				TLSHandshakeTimeout:   30 * time.Second,
				ResponseHeaderTimeout: 120 * time.Second,
				DisableKeepAlives:     true,
			},
		},
	}
}

type streamtapeServerResp struct {
	Status  int    `json:"status"`
	Msg     string `json:"msg"`
	Result  struct {
		URL string `json:"url"`
	} `json:"result"`
}

type streamtapeUploadResp struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Result struct {
		ID    string `json:"id"`
		URL   string `json:"url"`
		Embed string `json:"embed"`
	} `json:"result"`
}

// Upload uploads a file to Streamtape and returns the embed/view link
func (u *StreamtapeUploader) Upload(filePath string) (string, error) {
	uploadURL, err := u.getUploadURL()
	if err != nil {
		return "", fmt.Errorf("get upload URL: %w", err)
	}

	link, err := u.uploadFile(filePath, uploadURL)
	if err != nil {
		return "", fmt.Errorf("upload file: %w", err)
	}
	return link, nil
}

func (u *StreamtapeUploader) getUploadURL() (string, error) {
	url := fmt.Sprintf("https://api.streamtape.com/file/ul?login=%s&key=%s", u.login, u.key)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := u.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	var serverResp streamtapeServerResp
	if err := json.NewDecoder(resp.Body).Decode(&serverResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if serverResp.Status != 200 {
		return "", fmt.Errorf("API error %d: %s", serverResp.Status, serverResp.Msg)
	}
	if serverResp.Result.URL == "" {
		return "", fmt.Errorf("empty upload URL in response")
	}
	return serverResp.Result.URL, nil
}

func (u *StreamtapeUploader) uploadFile(filePath, uploadURL string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	pipeReader, pipeWriter := io.Pipe()
	writer := multipart.NewWriter(pipeWriter)
	errCh := make(chan error, 1)

	go func() {
		defer func() {
			writer.Close()
			pipeWriter.Close()
		}()
		part, err := writer.CreateFormFile("file", filepath.Base(filePath))
		if err != nil {
			errCh <- fmt.Errorf("create form file: %w", err)
			pipeWriter.CloseWithError(err)
			return
		}
		buf := make([]byte, 1024*1024)
		if _, err := io.CopyBuffer(part, file, buf); err != nil {
			errCh <- fmt.Errorf("copy file: %w", err)
			pipeWriter.CloseWithError(err)
			return
		}
		errCh <- nil
	}()

	req, err := http.NewRequest("POST", uploadURL, pipeReader)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := u.client.Do(req)
	if err != nil {
		select {
		case goroutineErr := <-errCh:
			if goroutineErr != nil {
				return "", fmt.Errorf("multipart write: %w (request: %v)", goroutineErr, err)
			}
		default:
		}
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload status %d: %s", resp.StatusCode, string(body))
	}

	var uploadResp streamtapeUploadResp
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return "", fmt.Errorf("decode upload response: %w", err)
	}
	if uploadResp.Status != 200 {
		return "", fmt.Errorf("upload API error %d: %s", uploadResp.Status, uploadResp.Msg)
	}
	if uploadResp.Result.ID == "" {
		return "", fmt.Errorf("empty file ID in upload response")
	}

	embedURL := uploadResp.Result.Embed
	if embedURL == "" {
		embedURL = fmt.Sprintf("https://streamtape.com/e/%s/", uploadResp.Result.ID)
	}
	return embedURL, nil
}
