package agent

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aryadhira/otogenius-agent/internal/llm"
	"github.com/aryadhira/otogenius-agent/internal/models"
)

type AgentExtractor struct {
	client     llm.LlmProvider
	masterdata []models.BrandModel
}

func NewAgentExtractor(client llm.LlmProvider, masterdata []models.BrandModel) Agent {
	return &AgentExtractor{
		client:     client,
		masterdata: masterdata,
	}
}

func (a *AgentExtractor) Run(prompt string) (any, error) {
	fmt.Println("\x1b[3m start parsing user requirement")
	systemMessage, err := getExtractorSystemPrompt(a.masterdata)
	if err != nil {
		return nil, err
	}

	messages := []models.Message{
		{Role: "system", Content: systemMessage},
		{Role: "user", Content: prompt},
	}

	response, err := a.client.ChatCompletions(messages, []models.Tool{})
	if err != nil {
		return nil, err
	}

	choice := response.Choices[0]
	aiResponse := choice.Message.Content.(string)

	fmt.Println("\x1b[3m parsed requirement", aiResponse)
	// sanitize AI response additional string
	sanitizedStr := strings.ReplaceAll(aiResponse, "```json", "")
	sanitizedStr = strings.ReplaceAll(sanitizedStr, "```", "")

	var result map[string]any

	// convert into map[string]any
	err = json.Unmarshal([]byte(sanitizedStr), &result)
	if err != nil {
		return nil, err
	}

	var filterPrompt strings.Builder
	for key, val := range result {
		str := fmt.Sprintf("%s : %v \n", key, val)
		filterPrompt.WriteString(str)
	}

	return filterPrompt.String(), nil
}

func getExtractorSystemPrompt(masterdata []models.BrandModel) (string, error) {
	var brand strings.Builder
	var model strings.Builder

	for _, each := range masterdata {
		brand.WriteString(fmt.Sprintf("%s,", each.BrandName))
		model.WriteString(fmt.Sprintf("%s,", each.ModelName))
	}

	sysMessageTemplate := `
		you are AI agent extractor assistant,
		your user is customer is someone who is looking for used car or seconhand car,
		your customer will describe their requirement of used car in natural language,
		your task is extract that requirement from their prompt.
		your output will be json in this schema :
		{
			brand : <string> optional (this will contains any brand that user mention, leave as empty string if user not mention any brand),
			model : <string> optional(this will contains any model that user mention, leave as empty string if user not mention any model),
			category : <string> mandatory (Enums: Sedan,SUV,MPV,Hatchback) you will clasify this value based on user prompt description judge by Body Style, Purpose, Seating Capacity, or Cargo Space,
			price : <integer> mandatory you will extract user budget or preference price, fill 0 if user not mention,
			production_year : <integer> optional you will extract any mentioned production_year of the car fill 0 if user not mention any year,
			transmission : <string> (Enums:Automatic,Manual) optional you will extract any mentioned transmission type from user,  leave as empty string if user not mention any transmission type,
		}
		
		this is example of minimal requirement json result  : 
		{
			"brand" : "",
			"model" : "",
			"category" : "Sedan",
			"price" : 200000000,
			"production_year" : 0,
			"transmission": ""
		}
		this is example of requirement with multiple brand and selected production year json result :
		{
			"brand" : "Toyota,Honda",
			"model" : "",
			"category" : "Sedan",
			"price" : 200000000,
			"production_year" : 2010,
			"transmission": ""
		}
		this is example of requirement with multiple brand & model and selected production year json result :
		{
			"brand" : "Toyota,Honda",
			"model" : "Civic,Corolla",
			"category" : "Sedan",
			"price" : 200000000,
			"production_year" : 2010,
			"transmission": ""
		}
		this is example of requirement with multiple brand & model and selected production year and selected transmission type json result :
		{
			"brand" : "Toyota,Honda",
			"model" : "Civic,Corolla",
			"category" : "Sedan",
			"price" : 200000000,
			"production_year" : 2010,
			"transmission": "Automatic"
		}

		for your additional reference :
		here list of brand name : %s.
		here list of model name : %s.
		
		here some word dictionary that related with Automatic Trasmission : metik, matic, triptonic, otomatis.

		PLEASE NOTE YOUR OUTPUT WILL BE ONLY JSON FOLLOWING THAT EXAMPLE ABOVE WITHOUT ANY ADDITIONAL WORD OR EXPLANATION, AND YOU CATEGORIZE THE REQUIREMENT ONLY BASED ON THE CUSTOMER PROMPT AND WIKI, NO NEED TO IMPROVIZE OR GIVE ASUMPTION.
	`

	return fmt.Sprintf(sysMessageTemplate, brand.String(), model.String()), nil
}

func (a *AgentExtractor) RunContinues(prompt string, messages []models.Message) (any, error){
	return nil, nil
}
