package script

import (
	"database/sql"
	"encoding/csv"
	"os"

	"github.com/aryadhira/otogenius-agent/internal/models"
	"github.com/google/uuid"
)

func CreateCarBrandModelMaster(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS master_brand_model(
		id TEXT PRIMARY KEY,
		brand_name TEXT,
		model_name TEXT,
		type_name TEXT 
	)`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	file, err := os.Open("docs/brand-model-type.csv")
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	vehicleTypeMap := make(map[string]string)
	vehicleTypeMap["1"] = "Sedan"
	vehicleTypeMap["2"] = "SUV"
	vehicleTypeMap["3"] = "MPV"
	vehicleTypeMap["4"] = "Hatchback"

	for i, each := range records {
		if i == 0 {
			continue
		}
		data := models.BrandModel{
			Id:        uuid.NewString(),
			BrandName: each[0],
			ModelName: each[1],
			TypeName:  vehicleTypeMap[each[2]],
		}

		err = SaveMasterData(db, data)
		if err != nil {
			return err
		}

	}

	return nil
}

func SaveMasterData(db *sql.DB, data models.BrandModel) error {
	query := `INSERT INTO master_brand_model (id, brand_name, model_name, type_name) 
			  VALUES ($1,$2,$3,$4)`

	_, err := db.Exec(query,
		data.Id,
		data.BrandName,
		data.ModelName,
		data.TypeName,
	)

	return err
}
