package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type OllamaClient struct {
	BaseURL string
}

func NewOllamaClient(baseURL string) *OllamaClient {
	return &OllamaClient{BaseURL: baseURL}
}

func (client *OllamaClient) GenerateResponse(prompt string) (string, error) {
	requestBody, err := json.Marshal(map[string]string{
		"model":  "llama2",
		"prompt": prompt,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", client.BaseURL+"/api/generate", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	webClient := &http.Client{}
	resp, err := webClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var responseBuffer strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		var responseChunk map[string]interface{}
		if err := json.Unmarshal([]byte(line), &responseChunk); err != nil {
			fmt.Printf("Error unmarshalling response chunk: %v\n", err)
			continue
		}
		if responseText, ok := responseChunk["response"].(string); ok {
			responseBuffer.WriteString(responseText)
		}
		if done, ok := responseChunk["done"].(bool); ok && done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	return responseBuffer.String(), nil
}
