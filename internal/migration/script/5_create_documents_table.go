package script

import "database/sql"

func CreateDocuments(db *sql.DB) error {
	query := `
		CREATE TABLE documents (
			id SERIAL PRIMARY KEY,
			content TEXT,
			embedding VECTOR(768)
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
