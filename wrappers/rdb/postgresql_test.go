package thundrardb

import (
	"database/sql"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setUpPostgresql(t *testing.T, dsn string) error {
	db, _ := sql.Open("postgres", dsn)
	err := db.Ping()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE table IF NOT EXISTS test(id int, type text)")
	assert.NoError(t, err)
	return nil
}

func TestPostgresqlIntegration(t *testing.T) {

	dsn := "postgres://user:userpass@localhost/db?sslmode=disable"
	err := setUpPostgresql(t, dsn)
	if err != nil {
		t.Skip()
	}
	s := newSuite(t, &pq.Driver{}, dsn, "postgresql")
	s.TestRdbIntegration(t, "SELECT * FROM test WHERE id = $1", "POSTGRESQL", 1)
}
