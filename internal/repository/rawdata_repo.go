package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/aryadhira/otogenius-agent/internal/models"
	"github.com/google/uuid"
)

type RawdataRepo interface {
	InsertRawData(ctx context.Context, rawdata *models.RawData) error
}

type RawDataImp struct {
	ctx context.Context
	db  *sql.DB
}

func NewRawData(ctx context.Context, db *sql.DB) RawdataRepo {
	return &RawDataImp{
		ctx: ctx,
		db:  db,
	}
}

func (r *RawDataImp) InsertRawData(ctx context.Context, rawdata *models.RawData) error {
	query := `INSERT INTO rawdata (id,brand,model,title,varian,fuel,transmission,image,price,source,scrape_date)
			  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`

	_, err := r.db.ExecContext(ctx, query,
		uuid.NewString(),
		rawdata.Brand,
		rawdata.Model,
		rawdata.Title,
		rawdata.Varian,
		rawdata.Fuel,
		rawdata.Transmission,
		rawdata.Image,
		rawdata.Price,
		rawdata.Source,
		time.Now(),
	)

	return err
}
