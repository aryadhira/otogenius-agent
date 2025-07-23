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

		strFilter := fmt.Sprintf("%s = '%s' ", key, val)
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
