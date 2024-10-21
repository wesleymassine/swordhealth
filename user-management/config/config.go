package config

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLConnection() (*sql.DB, error) {
    db, err := sql.Open(os.Getenv("SQL_DRIVER"), os.Getenv("MYSQL_URI"))
    if err != nil {
        return nil, err
    }
    return db, nil
}
