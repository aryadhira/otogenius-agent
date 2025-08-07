package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aryadhira/otogenius-agent/internal/agent"
	"github.com/aryadhira/otogenius-agent/internal/llm"
	"github.com/aryadhira/otogenius-agent/internal/migration"
	"github.com/aryadhira/otogenius-agent/internal/models"
	"github.com/aryadhira/otogenius-agent/internal/storages"
	"github.com/aryadhira/otogenius-agent/internal/tools"

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

	// llmModel := os.Getenv("OPENROUTER_MODEL")
	llmUrl := os.Getenv("LLM_URL")
	// apiKey := os.Getenv("OPENROUTER_API_KEY")

	// openrouter, err := llm.NewOpenRouter(llmUrl, llmModel, apiKey)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// ctx := context.Background()
	// masterdata := repository.NewBrandModel(ctx, db)

	temperatureStr := os.Getenv("TEMPERATURE")
	maxTokenStr := os.Getenv("MAX_TOKENS")

	temperature, _ := strconv.ParseFloat(temperatureStr, 64)
	maxToken, _ := strconv.Atoi(maxTokenStr)
	llamacpp := llm.NewLlamaCpp(llmUrl, temperature, maxToken, false)

	listTools := tools.RegisterTools()
	reader := bufio.NewReader(os.Stdin)

	// brandmodel, err := masterdata.GetAllBrandModel()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// extractor := agent.NewAgentExtractor(llamacpp, brandmodel)
	// recommendator := agent.NewAgentRecommendator(llamacpp, listTools)
	// advisor := agent.NewAgentAdvisor(llamacpp, brandmodel, listTools)
	// sysPrompt := agent.GetAdvisorSystemPrompt(listTools, brandmodel)

	// structured := agent.NewStructureOutput(llamacpp, listTools)

	otogenius := agent.NewOtogenius(llamacpp, listTools)
	messages := []models.Message{}

	fmt.Print("===================================================================================================\n")
	fmt.Println("--------Welcome to Otogenius Agent--------")
	fmt.Println("Please describe your used car requirement in natural languange your description can contains some key eg: brand, model, car production year, budget, transmission type")

	for {
		fmt.Print("===================================================================================================")
		fmt.Print("\nDescribe Your Requirement: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		// input = input + " /no_think"

		// messages = append(messages, models.Message{Role: "user", Content: input})

		// _, err := advisor.RunContinues(input, messages)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// _, err = advisor.Run(input)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// res, err := extractor.Run(input)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// time.Sleep(5 * time.Second)
		// _, err = recommendator.Run(res.(string))
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// _, err := structured.Run(input)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		history, err := otogenius.RunContinues(input, messages)
		if err != nil {
			log.Fatal(err)
		}
		messages = history.([]models.Message)
	}

}
