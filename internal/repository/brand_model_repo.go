package repository

import (
	"context"
	"database/sql"

	"github.com/aryadhira/otogenius-agent/internal/models"
)

type BrandModelRepo interface {
	GetAllBrandModel() ([]models.BrandModel, error)
}

type BrandModelImp struct {
	ctx context.Context
	db  *sql.DB
}

func NewBrandModel(ctx context.Context, db *sql.DB) BrandModelRepo {
	return &BrandModelImp{
		ctx: ctx,
		db:  db,
	}
}

func (b *BrandModelImp) GetAllBrandModel() ([]models.BrandModel, error) {
	query := `SELECT id, brand_name, model_name, type_name FROM master_brand_model`
	rows, err := b.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.BrandModel
	for rows.Next() {
		data := models.BrandModel{}
		err := rows.Scan(
			&data.Id,
			&data.BrandName,
			&data.ModelName,
			&data.TypeName,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, data)
	}

	return results, nil
}
