package cloudinary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/giakiet05/lkforum/internal/config"
)

const (
	apiBaseURL = "https://api.cloudinary.com/v1_1"
)

// Response is the structure of a successful Cloudinary API response.
type Response struct {
	SecureURL string `json:"secure_url"`
	PublicID  string `json:"public_id"`
}

// Upload sends a file to Cloudinary and returns the response.
func Upload(file multipart.File, fileHeader *multipart.FileHeader) (*Response, error) {
	cloudinaryCfg := config.Cfg.Cloudinary

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(part, file); err != nil {
		return nil, err
	}

	writer.WriteField("folder", cloudinaryCfg.UploadFolder)
	writer.WriteField("upload_preset", cloudinaryCfg.UploadPreset)

	writer.Close()

	url := fmt.Sprintf("%s/%s/image/upload", apiBaseURL, cloudinaryCfg.CloudName)
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(cloudinaryCfg.APIKey, cloudinaryCfg.APISecret)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("cloudinary upload error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var cloudRes Response
	if err := json.NewDecoder(resp.Body).Decode(&cloudRes); err != nil {
		return nil, err
	}

	return &cloudRes, nil
}

// Delete removes an image from Cloudinary using its public ID.
func Delete(publicID string) error {
	cloudinaryCfg := config.Cfg.Cloudinary

	url := fmt.Sprintf("%s/%s/image/destroy", apiBaseURL, cloudinaryCfg.CloudName)

	payload := map[string]string{
		"public_id": publicID,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.SetBasicAuth(cloudinaryCfg.APIKey, cloudinaryCfg.APISecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("cloudinary delete error (%d): %s", resp.StatusCode, data)
	}

	return nil
}
