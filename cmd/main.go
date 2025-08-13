package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/aryadhira/otogenius-agent/internal/agent"
	"github.com/aryadhira/otogenius-agent/internal/llm"
	"github.com/aryadhira/otogenius-agent/internal/migration"
	"github.com/aryadhira/otogenius-agent/internal/repository"
	"github.com/aryadhira/otogenius-agent/internal/services"
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

	carRepo := repository.NewCarRepo(context.Background(), db)
	embeddingRepo := repository.NewEmbeddingRepo(context.Background(), db)
	agent := agent.NewOtogenius(llamacpp, embed, embeddingRepo)

	otogenius := services.NewOtogeniusSvc(carRepo, agent)
	service := services.NewServiceHandler(otogenius)

	err = service.Start()
	if err != nil {
		log.Fatal(err)
	}

}
