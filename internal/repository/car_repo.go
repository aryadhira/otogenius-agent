package repository

import (
	"context"
	"database/sql"

	"github.com/aryadhira/otogenius-agent/internal/models"
)

type CarRepo interface {
	InsertCarData(ctx context.Context, rawdata models.CarInfo) error
	GetCarData(ctx context.Context, filter ...any) ([]models.CarInfo, error)
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

func (c *CarInfoImp) InsertCarData(ctx context.Context, rawdata models.CarInfo) error {

	return nil
}

func (c *CarInfoImp) GetCarData(ctx context.Context, filter ...any) ([]models.CarInfo, error) {
	return nil, nil
}
