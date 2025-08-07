package agent

import (
	"fmt"

	"github.com/aryadhira/otogenius-agent/internal/llm"
	"github.com/aryadhira/otogenius-agent/internal/models"
)

type StructureOutput struct {
	client llm.LlmProvider
	tools []models.Tool
}

func NewStructureOutput(client llm.LlmProvider, tools []models.Tool) Agent {
	return &StructureOutput{
		client: client,
		tools: tools,
	}
}

func getJsonSchema() map[string]any {
	schema := map[string]any{
		"type" : "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type": "string",
				"description": "full name of person",
			},
			"age": map[string]any{
				"type": "integer",
				"description": "The age of the person.",
			},
			"occupation": map[string]any{
				"type": "string",
				"description": "The primary occupation or field of work of the person.",
			},
			"notable_contribution": map[string]any{
				"type": "string",
				"description": "A significant contribution or achievement of the person.",
			},
		},
		"required" : []string{"name","occupation", "notable_contribution"},
	}
	return schema
}

func (s *StructureOutput)Run(prompt string) (any, error){
	sysMessage := "You are a helpful assistant that provides information about people in JSON format."

	message := []models.Message{
		{Role: "system", Content: sysMessage},
		{Role: "user", Content: prompt},
	}

	response, err := s.client.ChatCompletionsStructureOutput(message, []models.Tool{}, getJsonSchema())
	if err != nil {
		return nil, err
	}

	fmt.Println(response.Choices[0].Message.Content)

	return nil, nil

}

func (s *StructureOutput)RunContinues(prompt string, messages []models.Message) (any, error){
	return nil, nil
}