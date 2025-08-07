package agent

import (
	"encoding/json"
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

func getSystemPrompt(toolsDesc string) string {
	promptTemplate := `
		You are a specialized AI agent designed to assist customers in finding used cars. Your role is to identify a user's requirements from their query and trigger the appropriate function call.

		1. Available Tools
			You have access to one tool: get_car_catalog. Use this tool only when a user explicitly requests to search for a used car and provides all mandatory information.
		2. Tool Definitions
			%s
		3. Agents Behavior
			- Triggering the Function: You must call the find_used_car function whenever a user provides all the mandatory parameters (price and category). You should fill in any optional parameters you can identify.
			- Handling Missing Information: If the user's query is missing a mandatory parameter (price or category), do not attempt to call the function. Instead, you must politely and clearly ask the user for the missing information.
			- Constraint: You are an agent for searching used cars. If a user asks for something completely unrelated, you should respond with a friendly message stating that your capabilities are limited to finding cars. For example: "I can only help you find used cars. What kind of car are you looking for?"
	`
	return fmt.Sprintf(promptTemplate,toolsDesc)
}

func (s *Otogenius) Run(prompt string) (any, error) {
	return nil, nil
}

func (s *Otogenius) RunContinues(prompt string, messages []models.Message) (any, error) {
	fmt.Println("\x1b[3m\nstart retrieving recommendation based your requirement")
	if len(messages) == 0 {
		toolsJson, _ := json.Marshal(s.tools)

		messages = append(messages, models.Message{Role: "system", Content: getSystemPrompt(string(toolsJson))})
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
		usedCarList := []models.CarInfo{}
		listStr := toolResponse[len(toolResponse)-1].Content.(string)
		err = json.Unmarshal([]byte(listStr),&usedCarList)
		if err != nil {
			return nil, err
		}

		fmt.Println("\nOtogenius: ")
		for _,each := range usedCarList {
			fmt.Printf("Brand : %s \nModel : %s \nYear: %v \nPrice: %v \nTransmission: %s \nVarian: %s \n====================================\n", each.Brand, each.Model, each.ProductionYear, int(each.Price),each.Transmission, each.Varian)
		}
		fmt.Println("\x1b[0m")

		// fmt.Println("result :", toolResponse[len(toolResponse)-1].Content,"\x1b[0m")
	} else {
		messages = append(messages, models.Message{Role: "assistant", Content: message.Content})
		fmt.Printf("\nOtogenius: \n%s\n\x1b[0m", message.Content)
	}

	return messages, nil
}
