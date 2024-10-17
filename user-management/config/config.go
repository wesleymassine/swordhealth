package config

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func NewMySQLConnection() (*sql.DB, error) {
    db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/task")
    if err != nil {
        return nil, err
    }
    return db, nil
}
