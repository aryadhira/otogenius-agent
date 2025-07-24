package llamacpp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aryadhira/otogenius-agent/internal/models"
)

type LlamacppClient struct {
	url         string
	temperature float64
	maxTokens   int
	stream      bool
}

func NewLlamacppClient(url string, temperature float64, maxTokens int, stream bool) *LlamacppClient {
	return &LlamacppClient{
		url:         url,
		temperature: temperature,
		maxTokens:   maxTokens,
		stream:      stream,
	}
}

func (c *LlamacppClient) ChatCompletions(messages []models.Message, tools []models.Tool) (*models.LlmResponse, error) {
	reqBody := models.LlmRequest{
		Messages:       messages,
		Temperature:    c.temperature,
		MaxTokens:      c.maxTokens,
		Stream:         c.stream,
		Tools:          tools,
		ToolChoice:     "auto",
		ResponseFormat: map[string]string{"type": "text"},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

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
