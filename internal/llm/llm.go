package llm

import "github.com/aryadhira/otogenius-agent/internal/models"

type LlmProvider interface {
	ChatCompletions(messages []models.Message, tools []models.Tool) (*models.LlmResponse, error)
}
