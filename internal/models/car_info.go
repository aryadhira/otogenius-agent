package models

import "time"

type CarInfo struct {
	Id             string    `json:"id"`
	Brand          string    `json:"brand"`
	Model          string    `json:"model"`
	ProductionYear int       `json:"production_year"`
	Category       string    `json:"category"`
	Varian         string    `json:"varian"`
	Fuel           string    `json:"fuel"`
	Transmission   string    `json:"transmission"`
	ImageUrl       string    `json:"image_url"`
	Price          float64   `json:"price"`
	ScrapeDate     time.Time `json:"scrape_date"`
	ScrapeDateIn   int       `json:"scrape_dateint"`
}
