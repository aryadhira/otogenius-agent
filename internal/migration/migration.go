package migration

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/aryadhira/otogenius-agent/internal/migration/script"
	_ "github.com/lib/pq" // Import PostgreSQL driver
)

type DBMigration struct {
	db *sql.DB
}

func NewDBMigration(db *sql.DB) *DBMigration {
	return &DBMigration{
		db: db,
	}
}

func (m *DBMigration) StartMigration() error {
	dbVersion := m.checkCurrentDBVersion()
	log.Println("current db version:", dbVersion)

	if dbVersion == 0 {
		// create initial table for the first time
		err := m.createInitialTable()
		if err != nil {
			return err
		}

		return m.applyMigrations(1)

	}

	return m.applyMigrations(dbVersion)
}

func (m *DBMigration) createInitialTable() error {
	// Step 1: Create the "db_version" table if it does not exist
	createTableSQL := "CREATE TABLE IF NOT EXISTS db_version (version INTEGER NOT NULL);"
	_, err := m.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	log.Println("Table created successfully.")

	// Step 2: Insert the initial version (1) into the "version" column
	insertVersionSQL := "INSERT INTO db_version (version) VALUES (1);"
	_, err = m.db.Exec(insertVersionSQL)
	if err != nil {
		return fmt.Errorf("failed to insert initial version: %w", err)
	}
	log.Println("Initial version inserted successfully.")

	return nil
}

func (m *DBMigration) checkCurrentDBVersion() int {
	checkVersionSQL := "SELECT MAX(version) AS latest_version FROM db_version;"
	row := m.db.QueryRow(checkVersionSQL)
	latestVersion := 0
	err := row.Scan(&latestVersion)
	if err != nil {
		log.Println(err)
	}

	return latestVersion
}

func (m *DBMigration) applyMigrations(currentVersion int) error {
	for _, migration := range script.Migrations {
		if migration.Version > currentVersion {
			log.Printf("Applying migration version %d\n", migration.Version)
			if err := migration.Migrate(m.db); err != nil {
				return fmt.Errorf("migration %d failed: %w", migration.Version, err)
			}

			// update db_version table
			_, err := m.db.Exec("INSERT INTO db_version (version) VALUES ($1);", migration.Version)
			if err != nil {
				return fmt.Errorf("failed to update db_version: %w", err)
			}
			log.Println("Migration applied successfully")
		}
	}
	log.Println("DB Version up to date")

	return nil
}
