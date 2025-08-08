package repository

import (
	"context"
	"database/sql"

	"github.com/aryadhira/otogenius-agent/utils"
)

type EmbeddingRepo interface {
	SearchSimilarity(embedding []float32, limit int) ([]string, error)
}

type EmbeddingImp struct {
	ctx context.Context
	db  *sql.DB
}

func NewEmbeddingRepo(ctx context.Context, db *sql.DB) EmbeddingRepo {
	return &EmbeddingImp{
		ctx: ctx,
		db:  db,
	}
}

func (e *EmbeddingImp) SearchSimilarity(embedding []float32, limit int) ([]string, error) {
	query := `
		SELECT content
		FROM documents
		ORDER BY embedding <-> $1::vector
		LIMIT $2
	`

	vectorStr := utils.ConvertEmbeddingToString(embedding)

	rows, err := e.db.QueryContext(e.ctx, query, vectorStr, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			return nil, err
		}
		results = append(results, content)
	}

	return results, nil
}
