package storages

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", "postgres", "password", "localhost", "5432", "db_test")

	db, err := sql.Open("postgres", connString)
	assert.NotNil(t, db)
	assert.NoError(t, err)

	dbClone := db

	queryDrop := `DROP TABLE IF EXISTS test_table`
	_, err = dbClone.Query(queryDrop)
	assert.NoError(t, err)

	query := `CREATE TABLE IF NOT EXISTS test_table (
        column1 TEXT
    )`
	_, err = dbClone.Query(query)
	assert.NoError(t, err)

	query2 := "INSERT INTO test_table (column1) VALUES ($1)"
	stmt, err := dbClone.Prepare(query2)
	assert.NoError(t, err)
	_, err = stmt.Exec(1)
	assert.NoError(t, err)
}
