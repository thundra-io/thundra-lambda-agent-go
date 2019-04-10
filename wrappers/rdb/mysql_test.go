package thundrardb

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-sql-driver/mysql"
)

func setUpMysql(t *testing.T, dsn string) error {
	db, _ := sql.Open("mysql", dsn)
	err := db.Ping()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE table IF NOT EXISTS test(id int, type text)")
	assert.NoError(t, err)
	return nil
}

func TestMysqlIntegration(t *testing.T) {
	dsn := "user:userpass@tcp(localhost:3306)/db"
	err := setUpMysql(t, dsn)
	if err != nil {
		t.Skip()
	}
	s := newSuite(t, &mysql.MySQLDriver{}, dsn, "mysql")
	s.TestRdbIntegration(t, "SELECT * FROM test WHERE id = ?", "MYSQL", 1)
}
