package uploader

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// netuBase is the Netu.tv web API root.  The upload flow is documented at
	// https://netu.io/view_page.php?pid=10 (also mirrored in third-party clients
	// such as divseb/Netu.tv and Iliyass/node-netu).  Netu authenticates with a
	// user token (the "user_hash" you find on your profile page), NOT a separate
	// "api key" — the value configured as NETU_API_KEY is that user token.
	netuBase = "https://netu.tv"
)

// NetuUploader handles uploading files to Netu.tv (.tv / .ac).
//
// Flow (verified against Netu's published API + community clients):
//  1. GET  /plugins/cb_multiserver/api/get_upload_server.php?user_hash=<token>
//     -> upload_server, hash, time_hash, userid, key_hash, _remote
//  2. POST <upload_server> (multipart) with those tokens + Filedata + upload=1
//     -> success + file_name
//  3. POST /actions/file_uploader.php with insertVideo=yes, title, server,
//     user_hash, file_name -> video_link
//
// Netu does not return a poster/preview on upload, so thumbnails remain sourced
// from SeekStreaming / UPnShare only (see channel/pipeline.go stageSaveMetadata).
type NetuUploader struct {
	apiKey string // the Netu user token (user_hash)
	client *http.Client
}

// NewNetuUploader creates a new Netu uploader. An empty token disables the host
// (it is only registered when a token is configured).
func NewNetuUploader(apiKey string) *NetuUploader {
	return &NetuUploader{
		apiKey: strings.TrimSpace(apiKey),
		client: newDefaultClient(uploadClientTimeout),
	}
}

// netuServerResponse is step 1's payload.
type netuServerResponse struct {
	Successful   bool   `json:"successful"`
	UploadServer string `json:"upload_server"`
	Hash         string `json:"hash"`
	TimeHash     string `json:"time_hash"`
	UserID       string `json:"userid"`
	KeyHash      string `json:"key_hash"`
	Remote       string `json:"_remote"`
	Error        string `json:"error"`
}

// netuUploadResult is step 2's payload.
type netuUploadResult struct {
	Success  string `json:"success"`
	FileName string `json:"file_name"`
	Msg      string `json:"msg"`
}

// netuInsertResult is step 3's payload.
type netuInsertResult struct {
	VideoLink string `json:"video_link"`
	Success   string `json:"success"`
	Msg       string `json:"msg"`
}

// Upload uploads a file to Netu and returns the embed/download link.
func (u *NetuUploader) Upload(filePath string) (string, error) {
	return u.UploadWithProgress(filePath, nil)
}

// isNetuMaintenance reports whether the Netu API error indicates a server-side
// outage / maintenance window (e.g. the current "upload servers on upgrade, try
// again later" message). These are not retriable by us — Netu must finish its
// maintenance — so callers should treat them as permanent and skip Netu for now
// rather than burning every retry attempt on a guaranteed failure.
func isNetuMaintenance(msg string) bool {
	m := strings.ToLower(msg)
	return strings.Contains(m, "upload servers on upgrade") ||
		strings.Contains(m, "try again later") ||
		strings.Contains(m, "maintenance") ||
		strings.Contains(m, "temporarily unavailable") ||
		strings.Contains(m, "service unavailable")
}

// UploadWithProgress uploads a file to Netu and reports progress through fn.
func (u *NetuUploader) UploadWithProgress(filePath string, progress ProgressFunc) (string, error) {
	if u.apiKey == "" {
		return "", fmt.Errorf("Netu user token not configured")
	}

	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		if attempt > 1 {
			time.Sleep(uploadBackoff(attempt-2, lastErr))
		}
		link, err := u.uploadFile(filePath, progress)
		if err != nil {
			lastErr = fmt.Errorf("netu: %w", err)
			// Netu maintenance / outage: not retriable on our side. Mark as a
			// permanent error so the pipeline skips Netu immediately (without
			// burning the remaining retries) and proceeds with other hosts.
			if isNetuMaintenance(err.Error()) {
				return "", &permanentError{err: lastErr}
			}
			if isUploadRateLimited(err) || isQuotaExceeded(err.Error()) {
				// Rate limit / quota — back off and stop retrying this host.
				time.Sleep(uploadBackoff(attempt, err))
				return "", lastErr
			}
			if attempt < 3 {
				continue
			}
			return "", lastErr
		}
		return link, nil
	}
	return "", lastErr
}

func (u *NetuUploader) uploadFile(filePath string, progress ProgressFunc) (string, error) {
	// Step 1: get an upload server + tokens.
	var srv netuServerResponse
	if err := u.apiGet("/plugins/cb_multiserver/api/get_upload_server.php", "user_hash", u.apiKey, &srv); err != nil {
		return "", err
	}
	if srv.Error != "" {
		return "", fmt.Errorf("get_upload_server error: %s", srv.Error)
	}
	if srv.UploadServer == "" || srv.Hash == "" {
		return "", fmt.Errorf("incomplete upload_server response: %+v", srv)
	}

	// Step 2: upload the file to the upload server.
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	pipeReader, pipeWriter := io.Pipe()
	writer := multipart.NewWriter(pipeWriter)

	errChan := make(chan error, 1)
	go func() {
		defer pipeWriter.Close()
		for _, field := range []struct{ name, val string }{
			{"hash", srv.Hash},
			{"time_hash", srv.TimeHash},
			{"userid", srv.UserID},
			{"key_hash", srv.KeyHash},
			{"upload", "1"},
		} {
			if field.val == "" {
				continue
			}
			if err := writer.WriteField(field.name, field.val); err != nil {
				errChan <- fmt.Errorf("write field %s: %w", field.name, err)
				writer.Close()
				return
			}
		}
		part, err := writer.CreateFormFile("Filedata", filepath.Base(filePath))
		if err != nil {
			errChan <- fmt.Errorf("create form file: %w", err)
			writer.Close()
			return
		}
		fi, _ := file.Stat()
		var size int64
		if fi != nil {
			size = fi.Size()
		}
		progressFile := NewProgressReaderWithCallback(file, size, "Netu", progress)
		buf := make([]byte, 4*1024*1024)
		if _, err := io.CopyBuffer(part, progressFile, buf); err != nil {
			errChan <- fmt.Errorf("copy file: %w", err)
			writer.Close()
			return
		}
		if err := writer.Close(); err != nil {
			errChan <- fmt.Errorf("close writer: %w", err)
			return
		}
		errChan <- nil
	}()

	req, err := http.NewRequest("POST", srv.UploadServer, pipeReader)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := u.client.Do(req)
	if err != nil {
		pipeReader.CloseWithError(err)
		select {
		case <-errChan:
		case <-time.After(5 * time.Second):
		}
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	select {
	case err := <-errChan:
		if err != nil {
			return "", err
		}
	case <-time.After(30 * time.Second):
		return "", fmt.Errorf("timeout waiting for file copy to complete")
	}

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}
	var upResp netuUploadResult
	if err := json.Unmarshal(body, &upResp); err != nil {
		return "", fmt.Errorf("decode upload response: %w (body: %s)", err, string(body))
	}
	if upResp.Success == "" && upResp.FileName == "" {
		return "", fmt.Errorf("upload not accepted: %s (body: %s)", upResp.Msg, string(body))
	}
	if upResp.FileName == "" {
		return "", fmt.Errorf("upload returned empty file_name (body: %s)", string(body))
	}

	// Step 3: register the uploaded file to obtain its public video link.
	insertURL := netuBase + "/actions/file_uploader.php"
	form := url.Values{}
	form.Set("insertVideo", "yes")
	form.Set("title", filepath.Base(filePath))
	form.Set("server", srv.UploadServer)
	form.Set("user_hash", u.apiKey)
	form.Set("file_name", upResp.FileName)

	insReq, err := http.NewRequest("POST", insertURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("create insert request: %w", err)
	}
	insReq.Header.Set("User-Agent", defaultUserAgent)
	insReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	insResp, err := u.client.Do(insReq)
	if err != nil {
		return "", fmt.Errorf("insert request: %w", err)
	}
	defer insResp.Body.Close()
	insBody, _ := io.ReadAll(insResp.Body)
	if insResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("insert failed with status %d: %s", insResp.StatusCode, string(insBody))
	}
	var ins netuInsertResult
	if err := json.Unmarshal(insBody, &ins); err != nil {
		return "", fmt.Errorf("decode insert response: %w (body: %s)", err, string(insBody))
	}
	if ins.VideoLink == "" {
		return "", fmt.Errorf("insert returned empty video_link (body: %s)", string(insBody))
	}

	return ins.VideoLink, nil
}

func (u *NetuUploader) apiGet(path, keyParam, keyVal string, out interface{}) error {
	u2 := netuBase + path
	if strings.Contains(u2, "?") {
		u2 += "&"
	} else {
		u2 += "?"
	}
	u2 += keyParam + "=" + url.QueryEscape(keyVal)

	req, err := http.NewRequest("GET", u2, nil)
	if err != nil {
		return fmt.Errorf("netu: create request: %w", err)
	}
	req.Header.Set("User-Agent", defaultUserAgent)
	resp, err := u.client.Do(req)
	if err != nil {
		return fmt.Errorf("netu: do request: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("netu: HTTP %d: %s", resp.StatusCode, string(body))
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("netu: decode response: %w (body: %s)", err, string(body))
	}
	return nil
}
