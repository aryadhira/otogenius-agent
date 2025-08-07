package agent

import (
	"fmt"

	"github.com/aryadhira/otogenius-agent/internal/llm"
	"github.com/aryadhira/otogenius-agent/internal/models"
	"github.com/aryadhira/otogenius-agent/internal/tools"
)

type Otogenius struct {
	client llm.LlmProvider
	tools  []models.Tool
}

func NewOtogenius(client llm.LlmProvider, tools []models.Tool) Agent {
	return &Otogenius{
		client: client,
		tools:  tools,
	}
}

func getSystemPrompt() string {
	promptTemplate := `
		You are a specialized AI agent designed to assist customers in finding used cars. Your sole function is to identify a customer's requirements from their natural language query and use the provided get_car_catalog tool to search a database.

		Your process is as follows:

			1. Analyze the User's Request: Carefully read the customer's message to identify all relevant used car preferences.

			2. Extract Parameters: Map the identified preferences to the arguments of the get_car_catalog tool.

			3. Construct the Tool Call: Format a JSON object that strictly adheres to the tool's definition and parameter constraints.

			4. Output the Tool Call: Respond with only the completed JSON object. Do not add any conversational text, explanations, or extra characters.

		##Tool Definition##
		Tool Name: get_car_catalog

		Description: Searches the used car database for listings that match the specified criteria.

		Arguments:

		-	brand: A string. This is optional. If the user mentions multiple brands, they should be a comma-separated string (e.g., "Toyota,Honda,Mitsubishi"). If no brand is specified, the value must be an empty string "".

		-	model: A string. This is optional. If the user mentions multiple models, they should be a comma-separated string (e.g., "Civic,Corolla"). If no model is specified, the value must be an empty string "".

		-	price: An decimal. This is mandatory. The value must be a single number representing the maximum price the customer is willing to pay.

		-	category: A string. This is mandatory. The value must be one of the following exact enum values: "Sedan", "SUV", "Hatchback", or "MPV".

		-	production_year: An integer. This is optional. If specified, the value should be the production year (e.g., 2020). If no year is specified, the value must be a 0.

		-	transmission: A string. This is optional. The value must be one of the following exact enum values: "Manual", "Automatic", or an empty string "".

		Tool Call Format
		The output must be a single JSON object. The object must contain the tool name and an arguments object with all the keys defined in the get_car_catalog tool, populated with the extracted values.

		Example of a valid tool call :
		{
			"tool": "find_used_car",
			"arguments": {
				"brand": "toyota,honda",
				"model": "corolla",
				"price": 200000000,
				"category": "Sedan",
				"production_year": 2020,
				"transmission": "Automatic"
			}
		}
		
		If the customer's request is ambiguous, and you cannot determine a mandatory parameter (e.g., price or category), do not attempt to guess. Instead, you must respond with a simple, polite message asking for clarification.

		Example of a clarification response:
		I need to know the maximum price and car category you're looking for to find the best options for you.
	`
	return promptTemplate
}

func (s *Otogenius) Run(prompt string) (any, error) {
	return nil, nil
}

func (s *Otogenius) RunContinues(prompt string, messages []models.Message) (any, error) {
	if len(messages) == 0 {
		messages = append(messages, models.Message{Role: "system", Content: getSystemPrompt()})
	}

	messages = append(messages, models.Message{Role: "user", Content: prompt})

	response, err := s.client.ChatCompletions(messages, s.tools)
	if err != nil {
		return nil, err
	}

	choice := response.Choices[0]
	message := choice.Message

	if len(message.ToolCalls) > 0 {
		toolResponse, err := tools.ToolCalling(message)
		if err != nil {
			return nil, err
		}

		messages = append(messages, toolResponse...)
		fmt.Println("result :", toolResponse[len(toolResponse)-1].Content)
	} else {
		messages = append(messages, models.Message{Role: "assistant", Content: message.Content})
		fmt.Printf("\nOtogenius: \n%s\n\x1b[0m", message.Content)
	}

	return messages, nil
}
