package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type FileInfo struct {
	Name   string `json:"name"`
	IsOCR  bool   `json:"is_ocr"`
	DataID string `json:"data_id"`
}

type TaskRequestData struct {
	EnableFormula bool       `json:"enable_formula"`
	Language      string     `json:"language"`
	LayoutModel   string     `json:"layout_model"`
	EnableTable   bool       `json:"enable_table"`
	Files         []FileInfo `json:"files"`
}
type TaskResponseData struct {
	BatchID  string   `json:"batch_id"`
	FileURLs []string `json:"file_urls"`
}
type TaskResponse struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"` // Assuming there might be a message field for errors
	Data TaskResponseData `json:"data"`
}

const (
	BaseUrl = "https://mineru.net/api"
)

type Mineru struct {
	Token string
}

func NewMineru(token string) *Mineru {
	return &Mineru{token}
}

func (mineru *Mineru) CreateTask(taskData *TaskRequestData) (string, []string, error) {
	jsonData, err := json.Marshal(taskData)
	if err != nil {
		return "", nil, fmt.Errorf("failed to marshal task data: %w", err)
	}
	req, err := http.NewRequest("POST", BaseUrl+"/v4/file-urls/batch", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+mineru.Token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close() // Ensure body is closed
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var taskResp TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		return "", nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	if taskResp.Code != 0 {
		return "", nil, fmt.Errorf("API returned error code %d: %s", taskResp.Code, taskResp.Msg)
	}

	if len(taskResp.Data.FileURLs) == 0 {
		return "", nil, fmt.Errorf("API response successful but did not return file URLs")
	}
	fmt.Printf("Task creation successful. Batch ID: %s\n", taskResp.Data.BatchID)
	return taskResp.Data.BatchID, taskResp.Data.FileURLs, nil
}

func (mineru *Mineru) UploadFile(uploadUrl string, reader io.Reader) error {
	req, err := http.NewRequest("PUT", uploadUrl, reader)
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %w", err)
	}
	//	req.ContentLength = fileSize // Set content length

	client := &http.Client{Timeout: 30 * time.Second} // Adjust timeout as needed
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send upload request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("file upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	fmt.Printf("File uploaded successfully to %s\n", uploadUrl)
	return nil
}
