package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/aryadhira/otogenius-agent/internal/models"
)

type CarRepo interface {
	InsertCarData(carInfo models.CarInfo) error
	BulkInsertCarData(carsInfo []models.CarInfo) error
	GetCarData(filter map[string]any) ([]models.CarInfo, error)
}

type CarInfoImp struct {
	ctx context.Context
	db  *sql.DB
}

func NewCarRepo(ctx context.Context, db *sql.DB) CarRepo {
	return &CarInfoImp{
		ctx: ctx,
		db:  db,
	}
}

func (c *CarInfoImp) InsertCarData(carInfo models.CarInfo) error {
	query := `INSERT INTO car_info (id, brand, model, production_year, category, varian, fuel, transmission, image_url, price, scrape_date, scrape_dateint)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

	_, err := c.db.ExecContext(c.ctx, query,
		carInfo.Id,
		carInfo.Brand,
		carInfo.Model,
		carInfo.ProductionYear,
		carInfo.Category,
		carInfo.Varian,
		carInfo.Fuel,
		carInfo.Transmission,
		carInfo.ImageUrl,
		carInfo.Price,
		carInfo.ScrapeDate,
		carInfo.ScrapeDateInt,
	)

	return err
}

func (c *CarInfoImp) BulkInsertCarData(carsInfo []models.CarInfo) error {
	var err error

	const batchSize = 3000 // Maximum number of car data entries per batch insert
	columns := "id, brand, model, production_year, category, varian, fuel, transmission, image_url, price, scrape_date, scrape_dateint"

	for i := 0; i < len(carsInfo); i += batchSize {
		end := i + batchSize
		if end > len(carsInfo) {
			end = len(carsInfo)
		}
		currentBatch := carsInfo[i:end]

		// Slices to build the query and hold the values
		valuePlaceholders := make([]string, 0, len(currentBatch))
		valueArgs := make([]interface{}, 0, len(currentBatch)*12) // 12 columns per row

		// A counter for the positional parameters ($1, $2, ...)
		placeholderCounter := 1

		for _, car := range currentBatch {
			// Build the placeholders for a single row, e.g., "($1,$2,...,$12)"
			rowPlaceholders := make([]string, 12)
			for i := 0; i < 12; i++ {
				rowPlaceholders[i] = fmt.Sprintf("$%d", placeholderCounter)
				placeholderCounter++
			}
			valuePlaceholders = append(valuePlaceholders, fmt.Sprintf("(%s)", strings.Join(rowPlaceholders, ", ")))

			// Append the actual values to the args slice
			// The pq driver will handle the type conversions, including time.Time
			valueArgs = append(valueArgs,
				car.Id,
				car.Brand,
				car.Model,
				car.ProductionYear,
				car.Category,
				car.Varian,
				car.Fuel,
				car.Transmission,
				car.ImageUrl,
				car.Price,
				car.ScrapeDate,
				car.ScrapeDateInt,
			)
		}

		// Construct the final query string with placeholders
		query := fmt.Sprintf("INSERT INTO car_info (%s) VALUES %s",
			columns,
			strings.Join(valuePlaceholders, ", "))

		// Execute the query using the transaction and the valueArgs slice
		_, err = c.db.ExecContext(c.ctx, query, valueArgs...)
		if err != nil {
			return fmt.Errorf("failed to execute bulk insert query: %w", err)
		}
	}

	return err
}

func (c *CarInfoImp) GetCarData(filter map[string]any) ([]models.CarInfo, error) {
	var err error
	var results []models.CarInfo
	query := "SELECT id, brand, model, production_year, category, varian, fuel, transmission, image_url, price, scrape_date, scrape_dateint FROM car_info "
	finalQuery := parseFilter(query, filter)

	rows, err := c.db.QueryContext(c.ctx, finalQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		data := models.CarInfo{}
		err = rows.Scan(
			&data.Id,
			&data.Brand,
			&data.Model,
			&data.ProductionYear,
			&data.Category,
			&data.Varian,
			&data.Fuel,
			&data.Transmission,
			&data.ImageUrl,
			&data.Price,
			&data.ScrapeDate,
			&data.ScrapeDateInt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, data)
	}
	return results, nil
}

func parseFilter(query string, filter map[string]any) string {
	var sb strings.Builder
	sb.WriteString(query)

	filterCount := 0

	for key, val := range filter {
		if val == nil || val == "" {
			continue
		}

		strFilter := filterFormatter(key, val)
		if filterCount == 0 {
			sb.WriteString("WHERE ")
			sb.WriteString(strFilter)
		} else {
			andFilterString := fmt.Sprintf("AND %s", strFilter)
			sb.WriteString(andFilterString)
		}
		filterCount++
	}
	return sb.String()
}

func filterFormatter(key string, val any) string {
	filterFormat := ""
	switch key {
	case "brand", "model":
		strFilter := strings.Split(val.(string), ",")
		tempStr := []string{}
		for _, each := range strFilter {
			tempStr = append(tempStr, fmt.Sprintf("'%s' ", each))
		}
		filterFormat = fmt.Sprintf("%s IN (%s) ", key, strings.Join(tempStr, ", "))
	case "production_year", "price":
		filterFormat = fmt.Sprintf("%s <= %v ", key, val)
	default:
		filterFormat = fmt.Sprintf("%s = '%s' ", key, val)
	}

	return filterFormat
}
