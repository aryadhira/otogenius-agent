package scrapper

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aryadhira/otogenius-agent/internal/models"
	"github.com/aryadhira/otogenius-agent/internal/repository"
	"github.com/gocolly/colly/v2"
)

type OlxScrapper struct {
	ctx    context.Context
	db     repository.RawdataRepo
	master repository.BrandModelRepo
	colly  *colly.Collector
}

func NewOlxScrapper(ctx context.Context, db repository.RawdataRepo, master repository.BrandModelRepo, colly *colly.Collector) ScrapperInterface {
	return &OlxScrapper{
		ctx:    ctx,
		db:     db,
		master: master,
		colly:  colly,
	}
}

func (s *OlxScrapper) scrapAdsUrl(url string) []string {
	urls := []string{}
	s.colly.OnHTML("ul._21Jxw", func(h *colly.HTMLElement) {
		h.ForEach("li._3V_Ww a", func(i int, e *colly.HTMLElement) {
			href := e.Attr("href")
			urls = append(urls, href)
		})
	})

	s.colly.Visit(url)

	return urls
}

func (s *OlxScrapper) scrapAdsData(url, brand, model string) models.RawData {
	info := models.RawData{}
	info.Brand = brand
	info.Model = model
	info.ScrapeDate = time.Now()

	// Title {Brand Model (Year)}
	s.colly.OnHTML("h1._2iMMO", func(h *colly.HTMLElement) {
		info.Title = h.Text
	})

	// Varian
	s.colly.OnHTML("div.BxCeR", func(h *colly.HTMLElement) {
		info.Varian = h.Text
	})

	// index 0 = fuel
	// index 1 = transmission
	s.colly.OnHTML("div._1Im-S", func(h *colly.HTMLElement) {
		h.ForEach("h2._3rMkw", func(i int, e *colly.HTMLElement) {
			switch i {
			case 0:
				info.Fuel = e.Text
			case 1:
				info.Transmission = e.Text
			}
		})
	})

	// image
	s.colly.OnHTML("div._23Jeb", func(h *colly.HTMLElement) {
		src := h.ChildAttr("figure img", "src")
		info.Image = src
	})

	// price
	s.colly.OnHTML("div._1uqlc", func(h *colly.HTMLElement) {
		info.Price = h.Text
	})

	s.colly.Visit(url)

	return info
}

func (s *OlxScrapper) Run() error {
	start := time.Now()

	log.Println("Starting OLX scrapper...")

	baseUrl := os.Getenv("OLX_URL")

	master, err := s.master.GetAllBrandModel()
	if err != nil {
		return err
	}

	for _, each := range master {
		// preparing scrapper url
		modelName := strings.ReplaceAll(each.ModelName, " ", "-")
		queryString := strings.ToLower(fmt.Sprintf("%s-%s", each.BrandName, modelName))
		baseQueryUrl := "/jakarta-dki_g2000007/mobil-bekas_c198/q-"
		scrapUrlByModel := fmt.Sprintf("%s%s%s", baseUrl, baseQueryUrl, queryString)
		log.Println("scrapping : ", scrapUrlByModel)

		adsUrl := s.scrapAdsUrl(scrapUrlByModel)
		for _, url := range adsUrl {
			ads := baseUrl + url
			data := s.scrapAdsData(ads, each.BrandName, each.ModelName)
			err = s.db.InsertRawData(s.ctx, data)
			if err != nil {
				return err
			}
		}
		time.Sleep(1 * time.Second)
	}

	log.Println("Scrapping done at:", time.Since(start).Seconds())

	return nil
}
