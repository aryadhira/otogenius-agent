package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aryadhira/otogenius-agent/internal/llm"
	"github.com/aryadhira/otogenius-agent/internal/migration"
	"github.com/aryadhira/otogenius-agent/internal/models"
	"github.com/aryadhira/otogenius-agent/internal/repository"
	"github.com/aryadhira/otogenius-agent/internal/storages"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := storages.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	migs := migration.NewDBMigration(db)
	err = migs.StartMigration()
	if err != nil {
		log.Fatal(err)
	}

	llmUrl := os.Getenv("LLM_URL")
	embeddingUrl := os.Getenv("EMBEDDING_URL")
	temperatureStr := os.Getenv("TEMPERATURE")
	maxTokenStr := os.Getenv("MAX_TOKENS")

	temperature, _ := strconv.ParseFloat(temperatureStr, 64)
	maxToken, _ := strconv.Atoi(maxTokenStr)

	llamacpp := llm.NewLlamaCpp(llmUrl, temperature, maxToken, false)
	embed := llm.NewLlamaCpp(embeddingUrl, 0, 0, false)

	ctx := context.Background()
	embeddingRepo := repository.NewEmbeddingRepo(ctx, db)
	// masterdata := repository.NewBrandModel(ctx, db)

	// listTools := tools.RegisterTools()
	reader := bufio.NewReader(os.Stdin)

	// brandmodel, err := masterdata.GetAllBrandModel()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// otogenius := agent.NewOtogenius(llamacpp, listTools)
	promptTemplate := `
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
	// messages := []models.Message{}

	fmt.Print("===================================================================================================\n")
	fmt.Println("--------Welcome to Otogenius Agent--------")
	fmt.Println("Please describe your used car requirement in natural languange your description can contains some key eg: brand, model, car production year, budget, transmission type")

	for {
		fmt.Print("===================================================================================================")
		fmt.Print("\nDescribe Your Requirement: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		// input = input + " /no_think"
		embbeding, err := embed.GetEmbedding(input)
		if err != nil {
			log.Fatal(err)
		}

		messages := []models.Message{}

		similarDocs, err := embeddingRepo.SearchSimilarity(embbeding, 2)

		userMessage := fmt.Sprintf(promptTemplate, input, strings.Join(similarDocs, "\n"))
		messages = append(messages, models.Message{Role: "user", Content: userMessage})

		response, err := llamacpp.ChatCompletions(messages, []models.Tool{})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(response.Choices[0].Message.Content)

		// history, err := otogenius.RunContinues(input, messages)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// messages = history.([]models.Message)
	}

}
