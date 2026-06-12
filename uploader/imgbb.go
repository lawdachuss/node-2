package uploader

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const imgbbAPIURL = "https://api.imgbb.com/1/upload"
const imgbbDefaultKey = "5cddf83a7031b66400c06c622fdac2ad"

type imgbbResponse struct {
	Data struct {
		URL string `json:"url"`
	} `json:"data"`
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

type ImgBBUploader struct {
	apiKey string
	client *http.Client
}

func NewImgBBUploader() *ImgBBUploader {
	key := os.Getenv("IMGBB_API_KEY")
	if key == "" {
		key = imgbbDefaultKey
	}
	return &ImgBBUploader{
		apiKey: key,
		client: newNoProxyClient(60 * time.Second),
	}
}

func (u *ImgBBUploader) Upload(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("imgbb: read file: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(data)

	form := url.Values{
		"key":   {u.apiKey},
		"image": {encoded},
	}

	resp, err := u.client.PostForm(imgbbAPIURL, form)
	if err != nil {
		return "", fmt.Errorf("imgbb: post: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return "", fmt.Errorf("imgbb: read response: %w", err)
	}

	var result imgbbResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("imgbb: parse response: %w", err)
	}

	if result.Status != 200 {
		msg := result.Error
		if msg == "" {
			msg = string(body)
		}
		return "", fmt.Errorf("imgbb: error: %s", msg)
	}

	if result.Data.URL == "" {
		return "", fmt.Errorf("imgbb: empty image URL in response")
	}

	return result.Data.URL, nil
}
