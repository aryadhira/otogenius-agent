package agent

import (
	"encoding/json"
	"fmt"

	"github.com/aryadhira/otogenius-agent/internal/llamacpp"
	"github.com/aryadhira/otogenius-agent/internal/models"
	"github.com/aryadhira/otogenius-agent/internal/tools"
)

type AgentRecommendator struct {
	client *llamacpp.LlamacppClient
	tools  []models.Tool
}

func NewAgentRecommendator(client *llamacpp.LlamacppClient, tools []models.Tool) Agent {
	return &AgentRecommendator{
		client: client,
		tools:  tools,
	}
}

func (a *AgentRecommendator) Run(prompt string) (any, error) {
	fmt.Println("\x1b[3m start retrieve latest catalog list and find recommendation")
	systemMessage := getRecommendatorSystemPrompt(a.tools)
	userMessage := getUserPrompt(prompt)

	messages := []models.Message{
		{Role: "system", Content: systemMessage},
		{Role: "user", Content: userMessage},
	}

	err := a.executor(&messages)

	return nil, err
}

func (a *AgentRecommendator) executor(messages *[]models.Message) error {
	response, err := a.client.ChatCompletions(*messages, a.tools)
	if err != nil {
		return err
	}

	choice := response.Choices[0]
	aiResponse := choice.Message

	if len(aiResponse.ToolCalls) > 0 {
		history, err := tools.ToolCalling(aiResponse)
		if err != nil {
			return err
		}

		*messages = append(*messages, history...)
		// *messages = append(*messages, models.Message{Role: "user", Content: "if there is no data from the tools then tell user data with that criteria is not found, ask them to refine the prompt, but if there is data from tools summarize that catalog list into tabular data"})
		// *messages = append(*messages, models.Message{Role: "user", Content: "serve the result to user"})
		a.executor(messages)
	} else {
		fmt.Printf("Assistant: %s\n \x1b[0m", aiResponse.Content)
	}

	return nil
}

func getRecommendatorSystemPrompt(tools []models.Tool) string {
	promptTemplate := `
		you are AI Agent Recommendator, given access to these tools : 
		%s
		your task is providing latest catalog list of used car based on user request.
		please paid attention very carefully on car information that user provide, please pass exactly same for tool param don't assume any parameter.
		serve any data from tools as tabular data with this header : Brand, Model, Production Year, Category, Transmission, Price.
		show distinct data based on this combination : Brand, Model, Production Year, Category, Transmission, Average Price.
	`

	// if there is no data from the tools then tell user data with that criteria is not found, ask them to refine the prompt especially the mandatory field for tools.
	// but if there is data from tools summarize that catalog list into tabular data.

	toolsJson, _ := json.Marshal(tools)

	return fmt.Sprintf(promptTemplate, string(toolsJson))
}

func getUserPrompt(prompt string) string {
	promptTemplate := `
		please give me latest catalog of used car based on this information :
		%s
	`
	return fmt.Sprintf(promptTemplate, prompt)
}
