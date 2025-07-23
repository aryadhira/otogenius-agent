package transformation

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aryadhira/otogenius-agent/internal/models"
	"github.com/aryadhira/otogenius-agent/internal/repository"
	"github.com/google/uuid"
)

type Transformation interface {
	TransformCarInfoData() error
}

type TransformationImp struct {
	ctx        context.Context
	db         *sql.DB
	rawdata    repository.RawdataRepo
	carInfo    repository.CarRepo
	masterData repository.BrandModelRepo
}

func NewTransformation(ctx context.Context, db *sql.DB, rawdata repository.RawdataRepo, carInfo repository.CarRepo, masterData repository.BrandModelRepo) Transformation {
	return &TransformationImp{
		ctx:        ctx,
		db:         db,
		rawdata:    rawdata,
		carInfo:    carInfo,
		masterData: masterData,
	}
}

func (t *TransformationImp) TransformCarInfoData() error {
	start := time.Now()
	log.Println("Transformation started ....")
	// Get Raw Data
	rawdatas, err := t.rawdata.GetRawData(t.ctx)
	if err != nil {
		return err
	}

	masterData, err := t.masterData.GetAllBrandModel()
	if err != nil {
		return err
	}

	brandModelCategoryMap := make(map[string]string)
	for _, master := range masterData {
		key := fmt.Sprintf("%s-%s", master.BrandName, master.ModelName)
		brandModelCategoryMap[key] = master.TypeName
	}

	// Loop to normalize raw data into car info data
	log.Println("Start normalization")
	carInfos := []models.CarInfo{}
	for _, rawdata := range rawdatas {
		if rawdata.Title == "" {
			continue
		}

		carInfo := models.CarInfo{}
		carInfo.Id = uuid.NewString()
		carInfo.Brand = rawdata.Brand
		carInfo.Model = rawdata.Model
		carInfo.Varian = rawdata.Varian
		carInfo.Fuel = transformFuel(rawdata.Fuel)
		carInfo.Transmission = transformTransmission(rawdata.Transmission)
		carInfo.ImageUrl = rawdata.Image
		carInfo.ScrapeDate = rawdata.ScrapeDate

		categoryKey := fmt.Sprintf("%s-%s", rawdata.Brand, rawdata.Model)
		carInfo.Category = brandModelCategoryMap[categoryKey]

		scrapeDateStr := rawdata.ScrapeDate.Format("20060102")
		scrapeDateInt, _ := strconv.Atoi(scrapeDateStr)
		carInfo.ScrapeDateInt = int(scrapeDateInt)

		trimmedPrice := strings.TrimLeft(rawdata.Price, "Rp ")
		cleanPrice := strings.ReplaceAll(trimmedPrice, ".", "")
		priceFloat, _ := strconv.ParseFloat(cleanPrice, 64)
		carInfo.Price = priceFloat

		carInfo.ProductionYear, err = extractYearFromTitle(rawdata.Title)
		if err != nil {
			continue
		}

		carInfos = append(carInfos, carInfo)

	}
	log.Println("Normalization done")

	log.Println("Saving normalize data")
	// Insert Bulk Car Info
	err = t.carInfo.BulkInsertCarData(carInfos)
	if err != nil {
		return err
	}

	log.Printf("Transformation done on: %vs\n", time.Since(start).Seconds())

	return nil
}

func transformTransmission(transmission string) string {
	trns := ""
	if transmission == "Automatic" || transmission == "Manual" {
		trns = transmission
	}

	if strings.Contains(strings.ToLower(transmission), "otomatis") || strings.Contains(strings.ToLower(transmission), "option.automatic") || strings.Contains(strings.ToLower(transmission), "triptonic") {
		return "Automatic"
	}

	if strings.Contains(strings.ToLower(transmission), "option.manual") {
		return "Manual"
	}

	return trns
}

func transformFuel(fuel string) string {
	if strings.Contains(strings.ToLower(fuel), "option.petrol") {
		return "Bensin"
	}

	return fuel
}

func extractYearFromTitle(s string) (int, error) {
	re := regexp.MustCompile(`\((\d{4})\)`)
	matches := re.FindStringSubmatch(s)

	if len(matches) > 1 {
		// matches[0] is the entire match like "(2023)"
		// matches[1] is the content of the first capturing group, which is "2023"
		yearStr := matches[1]
		year, err := strconv.Atoi(yearStr) // Convert the string "2023" to an int 2023
		if err != nil {
			return 0, fmt.Errorf("failed to convert year string '%s' to int: %w", yearStr, err)
		}
		return year, nil
	}
	return 0, fmt.Errorf("no year found in parentheses in string: %s", s)
}
