package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aryadhira/otogenius-agent/internal/models"
)

type OpenRouter struct {
	url    string
	model  string
	apiKey string
}

type OpenRouterRequest struct {
	Model      string           `json:"model"`
	Messages   []models.Message `json:"messages"`
	Tools      []models.Tool    `json:"tools"`
	ToolChoice string           `json:"tool_choice"`
}

func NewOpenRouter(url string, model string, apiKey string) (LlmProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("apikey not found")
	}

	return &OpenRouter{
		url:    url,
		model:  model,
		apiKey: apiKey,
	}, nil
}

func (o *OpenRouter) ChatCompletions(messages []models.Message, tools []models.Tool) (*models.LlmResponse, error) {
	reqBody := OpenRouterRequest{
		Model:      o.model,
		Messages:   messages,
		Tools:      tools,
		ToolChoice: "auto",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, o.url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.apiKey))

	resp, err := clientRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Println(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server return error")
	}

	var LlmResponse models.LlmResponse

	err = json.Unmarshal(body, &LlmResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &LlmResponse, nil
}

func clientRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Retrying in 20 seconds...")
		time.Sleep(20 * time.Second)
		// IMPORTANT: Return the result of the recursive call
		return clientRequest(req)
	}

	return resp, nil
}

func (o *OpenRouter) ChatCompletionsStructureOutput(messages []models.Message, tools []models.Tool, jsonSchema map[string]any) (*models.LlmResponse, error) {
	return nil, nil
}

func (l *OpenRouter) GetEmbedding(text string) ([]float32, error) {
	return []float32{}, nil
}
