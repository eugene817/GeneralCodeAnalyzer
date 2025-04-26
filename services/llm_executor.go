package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GenerateRequest struct {
	Model  string `json:"model"`
	Stream bool   `json:"stream"`
	Prompt string `json:"prompt"`
}

type GenerateResponse struct {
	Model         string `json:"model"`
	Response      string `json:"response"`
	DoneReason    string `json:"done_reason"`
	TotalDuration int64  `json:"total_duration"`
}

func QueryOllama(prompt, model string) (string, error) {
	url := "http://localhost:11434/api/generate"

	// Formating the request body 
	reqBody := GenerateRequest{
		Model:  model,
		Stream: false,
		Prompt: prompt,
	}

	// Serialization of the request body
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// HTTP POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Reading the response body 
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// Unmarshalling the response body 
	var response GenerateResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// Checking the response status
	if response.DoneReason != "stop" {
		return "", fmt.Errorf("generation did not complete successfully: %s", response.DoneReason)
	}

	return response.Response, nil
}
