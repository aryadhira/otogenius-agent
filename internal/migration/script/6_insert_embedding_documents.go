package script

import (
	"database/sql"
	"os"

	"github.com/aryadhira/otogenius-agent/internal/llm"
	"github.com/aryadhira/otogenius-agent/utils"
)

func InsertEmbeddingDocuments(db *sql.DB) error {
	content := getWikiDocuments()
	query := `
		INSERT INTO documents (content, embedding)
		VALUES ($1, $2)
	`
	embeddingUrl := os.Getenv("EMBEDDING_URL")
	llamacpp := llm.NewLlamaCpp(embeddingUrl, 0, 0, false)

	chunks := utils.SplitIntoChunks(content, 15)

	for _, chunk := range chunks {
		embedding, err := llamacpp.GetEmbedding(chunk)
		if err != nil {
			return err
		}

		vectorStr := utils.ConvertEmbeddingToString(embedding)

		_, err = db.Exec(query, chunk, vectorStr)

		if err != nil {
			return err
		}
	}

	return nil
}

func getWikiDocuments() string {
	return `
		Question : What is best car category for family trip?
		Answer : MPV

		Question : What is best car category for camping, hiking, offroad?
		Answer : SUV

		Question : What is best car category that practical and easy to park?
		Answer: Hatchback

		Question : What is best car category that sporty and stylish?
		Answer: Sedan

		Question : What is best car category for business operational vehicle?
		Answer: MPV

		Question : What is best car category for school?
		Answer: Hatchback

		Question : What is best car category that can accomodate a lot passengers?
		Answer: MPV

		Question : What is most fuel efficient car category?
		Answer: Hatchback

		Question : What is best transmission type to passing traffic jam?
		Answer: Automatic
	`
}
