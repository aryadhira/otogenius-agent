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
	// rawdata := repository.NewRawData(ctx, db)
	masterdata := repository.NewBrandModel(ctx, db)
	// c := colly.NewCollector()
	// scrp := scrapper.NewOlxScrapper(ctx, rawdata, masterdata, c)

	// err = scrp.Run()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// carInfo := repository.NewCarRepo(ctx, db)
	// transform := transformation.NewTransformation(ctx, db, rawdata, carInfo, masterdata)
	// err = transform.TransformCarInfoData()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// filter := make(map[string]any)
	// filter["brand"] = "Toyota,Honda,Mitsubishi"
	// filter["model"] = "Corolla,Civic"
	// filter["category"] = "Sedan"
	// filter["price"] = 185000000
	// filter["production_year"] = 2015

	// res, err := carInfo.GetCarData(filter)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, each := range res {
	// 	data := fmt.Sprintf("%s %s %s %v %s %v", each.Brand, each.Model, each.Category, each.ProductionYear, each.Varian, int(each.Price))
	// 	log.Println(data)
	// }

	temperatureStr := os.Getenv("TEMPERATURE")
	maxTokenStr := os.Getenv("MAX_TOKENS")

	temperature, _ := strconv.ParseFloat(temperatureStr, 64)
	maxToken, _ := strconv.Atoi(maxTokenStr)
	url := os.Getenv("LLM_URL")

	client := llamacpp.NewLlamacppClient(url, temperature, maxToken, false)
	// listTools := tools.RegisterTools()
	reader := bufio.NewReader(os.Stdin)

	brandmodel, err := masterdata.GetAllBrandModel()
	if err != nil {
		log.Fatal(err)
	}

	extractor := agent.NewAgentExtractor(client, brandmodel)

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

		fmt.Println(res)
	}

}
