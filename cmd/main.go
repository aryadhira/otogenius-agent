package main

import (
	"context"
	"log"

	"github.com/aryadhira/otogenius-agent/internal/repository"
	"github.com/aryadhira/otogenius-agent/internal/scrapper"
	"github.com/aryadhira/otogenius-agent/internal/storages"
	"github.com/gocolly/colly/v2"
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

	ctx := context.Background()
	rawdata := repository.NewRawData(ctx, db)

	c := colly.NewCollector()
	olx := scrapper.NewOlxScrapper(ctx, rawdata, c)
	err = olx.Run()
	if err != nil {
		log.Fatal(err)
	}
}
