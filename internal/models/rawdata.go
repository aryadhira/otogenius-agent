package models

import "time"

type RawData struct {
	Id           string    `json:"id"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	Title        string    `json:"title"`
	Varian       string    `json:"varian"`
	Fuel         string    `json:"fuel"`
	Transmission string    `json:"transmission"`
	Image        string    `json:"image"`
	Price        string    `json:"price"`
	Source       string    `json:"source"`
	ScrapeDate   time.Time `json:"scrape_date"`
}
