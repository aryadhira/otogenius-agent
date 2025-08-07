package agent

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aryadhira/otogenius-agent/internal/llm"
	"github.com/aryadhira/otogenius-agent/internal/models"
	"github.com/aryadhira/otogenius-agent/internal/tools"
)

type AgentAdvisor struct {
	client     llm.LlmProvider
	masterdata []models.BrandModel
	tools      []models.Tool
}

func NewAgentAdvisor(client llm.LlmProvider, masterdata []models.BrandModel, tools []models.Tool) Agent {
	return &AgentAdvisor{
		client:     client,
		masterdata: masterdata,
		tools:      tools,
	}
}

func GetAdvisorSystemPrompt(tools []models.Tool, masterdata []models.BrandModel) string {
	var brand strings.Builder
	var model strings.Builder

	toolsJson, _ := json.Marshal(tools)

	for _, each := range masterdata {
		brand.WriteString(fmt.Sprintf("%s,", each.BrandName))
		model.WriteString(fmt.Sprintf("%s,", each.ModelName))
	}

	promptTemplate := `
		you are helpful Advisor AI agent, given access to these tools to access internal catalog car database:
		%s
		your user is customer is someone who is looking for used car or seconhand car,
		your customer will describe their requirement of used car in natural language,
		your task is giving customer car recommendation based on their requirement from our internal catalog database,
		since car category and price is mandatory for accessing internal catalog database please remember this before:
		- category -> you will retrieve this based on user requirement if you are not sure about the requirement please ask again user to refine their requirement
		- price -> if user not explicitly mention about budget or price please ask user to add the budget
		- brand -> if user not specify any brand no need to do assumption on that
		- model -> if user not specify any model no need to do assumption on that

		Important Notes :
		DON'T USE YOUR ASUMPTION for production_year.

		here reference list of brand on our database :
		%s
		here reference list of model on our database : 
		%s
		here some word dictionary that related with Automatic Trasmission : metik, matic, triptonic, otomatis.
		DON'T USE YOUR ASUMPTION TO ANSWER YOU WILL ONLY ANSWER BASED ON THE TOOL RESULT, PLEASE REMEMBER IT!

		Result format :
		- if you already get the result from tool, serve that list of used car data into tabular data with header : Brand, Model, Production Year, Transmission, Fuel, Varian, Price
		- if you get empty result from tool, just said "Can't find any used car in that specific requirement".
	`

	return fmt.Sprintf(promptTemplate, string(toolsJson), brand, model)

}

func GetAdvisorSystemPromptStructured(tools []models.Tool, masterdata []models.BrandModel) string {
	var brand strings.Builder
	var model strings.Builder

	toolsJson, _ := json.Marshal(tools)

	for _, each := range masterdata {
		brand.WriteString(fmt.Sprintf("%s,", each.BrandName))
		model.WriteString(fmt.Sprintf("%s,", each.ModelName))
	}

	promptTemplate := `
		you are helpful Advisor AI agent, given access to these tools to access internal catalog car database:
		%s
		your user is customer is someone who is looking for used car or seconhand car,
		your customer will describe their requirement of used car in natural language,
		your task is giving customer car recommendation based on their requirement from our internal catalog database,
		since car category and price is mandatory for accessing internal catalog database please remember this before accessing tools:
		- category -> you will retrieve this based on user requirement if you are not sure about the requirement please ask again user to refine their requirement
		- price -> if user not explicitly mention about budget or price please ask user to add the budget
		- brand -> if user not specify any brand no need to do assumption on that
		- model -> if user not specify any model no need to do assumption on that

		Important Notes :
		DON'T USE YOUR ASUMPTION for production_year.

		here reference list of brand on our database :
		%s
		here reference list of model on our database : 
		%s
		here some word dictionary that related with Automatic Trasmission : metik, matic, triptonic, otomatis.
		DON'T USE YOUR ASUMPTION TO ANSWER YOU WILL ONLY ANSWER BASED ON THE TOOL RESULT, PLEASE REMEMBER IT!

		Your Result format is JSON array must be follow this example:
		[
			{"brand" : "Toyota","model" : "Corolla","category" : "Sedan","price" : 200000000,"production_year" : 2013,"transmission": "Automatic"},
			{"brand" : "Honda","model" : "Civic","category" : "Sedan","price" : 200000000,"production_year" : 2014,"transmission": "Manual"},
		]
	`

	return fmt.Sprintf(promptTemplate, string(toolsJson), brand, model)

}

func (a *AgentAdvisor) advisorExecutor(messages *[]models.Message) error {
	response, err := a.client.ChatCompletions(*messages, a.tools)
	if err != nil {
		return err
	}

	choice := response.Choices[0]
	aiResponse := choice.Message

	if len(aiResponse.ToolCalls) > 0 {
		history, err := tools.ToolCalling(aiResponse)
		if err != nil {
			if strings.Contains(err.Error(), "can't parse parameter") {
				a.advisorExecutor(messages)
			}
			return err
		}

		*messages = append(*messages, history...)
		a.advisorExecutor(messages)
	} else {
		*messages = append(*messages, models.Message{Role: "assistant", Content: aiResponse.Content})
		fmt.Printf("\nOtogenius: \n%s\n\x1b[0m", aiResponse.Content)
	}

	return nil
}

func getUsedCarSchema() map[string]any {
	schema := map[string]any{
		"type": "array",
		"items": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"brand": map[string]any{
					"type":        "string",
					"description": "the car brand",
				},
				"model": map[string]any{
					"type":        "integer",
					"description": "the car brand model",
				},
				"category": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"Sedan", "SUV", "MPV", "Hatchback"},
					"description": "The car category",
				},
				"transmission": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"", "Automatic", "Manual"},
					"description": "The car transmission type",
				},
				"production_year": map[string]interface{}{
					"type":        "integer",
					"description": "Car production year",
				},
				"price": map[string]interface{}{
					"type":        "number",
					"description": "Car price",
				},
			},
			"required": []string{"brand", "model", "category", "transmission", "production_year", "price"},
		},
	}
	return schema
}

func (a *AgentAdvisor) advisorExecutorStructured(messages *[]models.Message) error {
	response, err := a.client.ChatCompletionsStructureOutput(*messages, a.tools, getUsedCarSchema())
	if err != nil {
		return err
	}

	choice := response.Choices[0]
	aiResponse := choice.Message

	if len(aiResponse.ToolCalls) > 0 {
		history, err := tools.ToolCalling(aiResponse)
		if err != nil {
			if strings.Contains(err.Error(), "can't parse parameter") {
				a.advisorExecutorStructured(messages)
			}
			return err
		}

		*messages = append(*messages, history...)
		a.advisorExecutorStructured(messages)
	} else {
		fmt.Println("assistant")
		*messages = append(*messages, models.Message{Role: "assistant", Content: aiResponse.Content})
		fmt.Println(aiResponse.Content)
		cars := []models.CarInfo{}
		res := aiResponse.Content.(string)
		err = json.Unmarshal([]byte(res), &cars)
		if err != nil {
			return err
		}

		fmt.Println(cars)
	}

	return nil
}

func (a *AgentAdvisor) Run(prompt string) (any, error) {

	fmt.Println("\x1b[3m\n\nstart retrieving recommendation based your requirement")
	systemMessage := GetAdvisorSystemPromptStructured(a.tools, a.masterdata)

	messages := []models.Message{
		{Role: "system", Content: systemMessage},
		{Role: "user", Content: prompt},
	}

	err := a.advisorExecutorStructured(&messages)

	return nil, err
}

func (a *AgentAdvisor) RunContinues(prompt string, messages []models.Message) (any, error) {
	fmt.Println("\x1b[3m\n\nstart retrieving recommendation based your requirement")
	err := a.advisorExecutor(&messages)

	return nil, err
}
