package script

import "database/sql"

func CreateRawDataTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS rawdata(
				id text primary key,
				brand text,
				model text,
				title text,
				varian text,
				fuel text,
				transmission text,
				image text,
				price text,
				source text,
				scrape_date timestamp
			)`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
