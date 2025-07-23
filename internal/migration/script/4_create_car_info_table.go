package script

import "database/sql"

func CreateCarInfo(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS car_info(
				id text primary key,
				brand text,
				model text,
				production_year integer,
				category text,
				varian text,
				fuel text,
				transmission text,
				image_url text,
				price numeric,
				scrape_date timestamp,
				scrape_dateint integer
			)`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
