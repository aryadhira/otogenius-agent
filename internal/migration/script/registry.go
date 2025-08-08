package script

import "database/sql"

type MigrationScript struct {
	Version int
	Migrate func(db *sql.DB) error
}

var Migrations = []MigrationScript{
	{Version: 2, Migrate: CreateCarBrandModelMaster},
	{Version: 3, Migrate: CreateRawDataTable},
	{Version: 4, Migrate: CreateCarInfo},
	{Version: 5, Migrate: CreateDocuments},
	{Version: 6, Migrate: InsertEmbeddingDocuments},
}
