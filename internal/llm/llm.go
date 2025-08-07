package llm

import "github.com/aryadhira/otogenius-agent/internal/models"

type LlmProvider interface {
	ChatCompletions(messages []models.Message, tools []models.Tool) (*models.LlmResponse, error)
	ChatCompletionsStructureOutput(messages []models.Message, tools []models.Tool, jsonSchema map[string]any) (*models.LlmResponse, error)
}
