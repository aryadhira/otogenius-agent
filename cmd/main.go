package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aryadhira/otogenius-agent/internal/agent"
	"github.com/aryadhira/otogenius-agent/internal/llamacpp"
	"github.com/aryadhira/otogenius-agent/internal/migration"
	"github.com/aryadhira/otogenius-agent/internal/repository"
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

	ctx := context.Background()
	masterdata := repository.NewBrandModel(ctx, db)

	// carInfo := repository.NewCarRepo(ctx, db)

	temperatureStr := os.Getenv("TEMPERATURE")
	maxTokenStr := os.Getenv("MAX_TOKENS")

	temperature, _ := strconv.ParseFloat(temperatureStr, 64)
	maxToken, _ := strconv.Atoi(maxTokenStr)
	url := os.Getenv("LLM_URL")

	client := llamacpp.NewLlamacppClient(url, temperature, maxToken, false)
	listTools := tools.RegisterTools()
	reader := bufio.NewReader(os.Stdin)

	brandmodel, err := masterdata.GetAllBrandModel()
	if err != nil {
		log.Fatal(err)
	}

	extractor := agent.NewAgentExtractor(client, brandmodel)
	recommendator := agent.NewAgentRecommendator(client, listTools)

	fmt.Print("===================================================================================================\n")
	fmt.Println("--------Welcome to Otogenius Agent--------")
	fmt.Println("Please describe your used car requirement in natural languange your description can contains some key eg: brand, model, car production year, budget, transmission type")

	for {
		fmt.Print("===================================================================================================")
		fmt.Print("\nDescribe Your Requirement: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		res, err := extractor.Run(input)
		if err != nil {
			log.Fatal(err)
		}

		_, err = recommendator.Run(res.(string))
		if err != nil {
			log.Fatal(err)
		}
	}

}
