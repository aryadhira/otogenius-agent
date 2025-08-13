package agent

import (
	"fmt"
	"strings"

	"github.com/aryadhira/otogenius-agent/internal/llm"
	"github.com/aryadhira/otogenius-agent/internal/models"
	"github.com/aryadhira/otogenius-agent/internal/repository"
)

type Otogenius struct {
	llm           llm.LlmProvider
	embed         llm.LlmProvider
	embeddingRepo repository.EmbeddingRepo
}

func NewOtogenius(llm, embed llm.LlmProvider, repo repository.EmbeddingRepo) Agent {
	return &Otogenius{
		llm:           llm,
		embed:         embed,
		embeddingRepo: repo,
	}
}

func getPromptTemplate() string {
	return `
		please extract this user used car preferences query : 
		%s
		
		your output will be json in this schema :
		{
			brand : <string> optional (this will contains any brand that user mention, leave as empty string if user not mention any brand),
			model : <string> optional(this will contains any model that user mention, leave as empty string if user not mention any model),
			category : <string> mandatory (Enums: Sedan,SUV,MPV,Hatchback) you will clasify this value based on this context %s,
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
	`
}

func (s *Otogenius) Run(prompt string) (any, error) {
	embbeding, err := s.embed.GetEmbedding(prompt)
	if err != nil {
		return nil, err
	}

	similarDocs, err := s.embeddingRepo.SearchSimilarity(embbeding, 2)
	if err != nil {
		return nil, err
	}

	messages := []models.Message{}
	userMessage := fmt.Sprintf(getPromptTemplate(), prompt, strings.Join(similarDocs, "\n"))
	messages = append(messages, models.Message{Role: "user", Content: userMessage})

	response, err := s.llm.ChatCompletions(messages, []models.Tool{})
	if err != nil {
		return nil, err
	}

	responseStr := response.Choices[0].Message.Content.(string)

	return responseStr, nil
}

func (s *Otogenius) RunContinues(prompt string, messages []models.Message) (any, error) {

	return nil, nil
}
