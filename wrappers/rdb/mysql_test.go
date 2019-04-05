package thundrardb

import (
	"database/sql"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func setUpMysql(t *testing.T, dsn string) error {
	db, _ := sql.Open("mysql", dsn)
	err := db.Ping()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE table IF NOT EXISTS test(id int, type text)")
	require.NoError(t, err)
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
